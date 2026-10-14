import { preloadPostCtaAssets, resolveBundleBaseUrl } from './assets.js';
import {
  BACKGROUND_WAVE_LIMITS,
  BACKGROUND_WAVE_STYLES,
  buildBackgroundPanoramaContext,
  extendedBackgroundPalette,
  paletteWithBackgroundHueShift,
  parseBackgroundWaveConfig,
  waveColorStrength,
  waveVisualStrength,
} from './background-panorama.js';
import { buildMotifStripContext } from './motif-strip.js';
import {
  installVisionStripController,
  restoreVisionStripScroll,
  showVisionStripPlaceholder,
  VISION_STRIP_PANEL_ID,
} from './connector-vision-export.js';
import { downloadCanvasWebp, exportFilename, variantIdFromIndex } from './export.js';
import { mixHex } from './background.js';
import { CAROUSEL_PALETTES, matchPaletteId, paletteBaseColorEntries, wavePaletteColorEntries, resolvePaletteId } from './palettes.js';
import {
  parseQrSizePercent,
  previewPixelDimensions,
  renderSlideToCanvas,
  compactSectionAlignment,
  normalizeSectionAlignment,
  sectionAlignmentFor,
} from './renderer.js';
import {
  buildThemeRenderOverrides,
  DESIGN_CANVAS_SIZE,
  loadDeckFonts,
  mergeTheme,
  parseAspectRatio,
  resolveColor,
  rebuildBackgroundGradient,
  deckPaletteFromTheme,
  deckWavePaletteFromTheme,
  isWavePaletteLinked,
  normalizeWavePalette,
  panoramaPaletteFromTextPalette,
  wavePaletteFromTheme,
  SLIDE_SIZE_MIN,
  SLIDE_SIZE_MAX,
} from './theme.js';
import {
  readIncludedSlideNumbers,
  readStripVariantIndices,
  resetIncludedSlideNumbers,
  resetStripVariantIndices,
  stripVariantIndexFor,
  writeIncludedSlideNumbers,
  writeStripVariantIndices,
} from './studio-slide-inclusion.js';
import {
  applyStudioThemeState,
  writeStudioThemeState,
} from './studio-theme-persistence.js';
import {
  applyStudioScrollState,
  bindStudioScrollPersistence,
  readStudioScrollState,
} from './studio-scroll-persistence.js';
import {
  captureFileDeckSnapshot,
  captureFileStripSnapshot,
  renderStudioStatePanel,
} from './studio-state-panel.js';
import { copyText, createCopyButton } from './copy-button.js';
import { createDownloadButton } from './download-button.js';
import { saveCarouselDeckToSource } from './studio-save.js';
import { CAROUSEL_SLIDE_WIDTH_PX, PANORAMA_SLIDE_WIDTH_PX, isCarouselCtaRole } from './slide-constants.js';

/** @type {{ paletteCards?: HTMLButtonElement[], gradientCards?: HTMLButtonElement[], waveModifierPanel?: HTMLElement, waveGeometryFields?: HTMLElement, waveInputs?: HTMLInputElement[], hueInputs?: HTMLInputElement[], colorPickers?: HTMLInputElement[], waveColorPickers?: HTMLInputElement[], wavePaletteSection?: HTMLElement, wavePaletteLinkToggle?: HTMLInputElement, lineHeightInputs?: HTMLInputElement[], marginInputs?: HTMLInputElement[], ctaInputs?: HTMLInputElement[], motifOffsetInputs?: HTMLInputElement[] }} */
const colorControlRefs = {};

/** @type {readonly { key: 'intensity' | 'color' | 'variety' | 'blur' | 'radius' | 'phase', label: string }[]} */
const WAVE_SCALAR_CONTROL_KEYS = [
  { key: 'intensity', label: 'Intensity' },
  { key: 'color', label: 'Color' },
  { key: 'variety', label: 'Color variety' },
  { key: 'blur', label: 'Blur' },
  { key: 'radius', label: 'Radius' },
  { key: 'phase', label: 'Phase' },
];

/** Deck margin inset as % of the 1080 design canvas (matches theme.js clamp). */
const MARGIN_PERCENT_MIN = 0;
const MARGIN_PERCENT_MAX = 16;

/** Default `post_cta` layout sizes at 1080 canvas (match renderer.js). */
const CTA_LAYOUT_DEFAULTS = {
  featuredMaxHeight: 200,
  qrSize: 100,
  brandMaxHeight: 180,
};

/** @type {Record<keyof typeof CTA_LAYOUT_DEFAULTS, { min: number, max: number, step: number, label: string }>} */
const CTA_LAYOUT_LIMITS = {
  featuredMaxHeight: { min: 100, max: 420, step: 5, label: 'Image max H' },
  qrSize: { min: 10, max: 100, step: 5, label: 'QR size' },
  brandMaxHeight: { min: 48, max: 320, step: 4, label: 'Logo max H' },
};

/** @type {readonly { key: string, label: string }[]} */
const LINE_HEIGHT_CONTROL_KEYS = [
  { key: 'normal', label: 'Body' },
  { key: 'punch', label: 'Punch' },
  { key: 'header', label: 'Header' },
  { key: 'footer', label: 'Footer' },
];

/** Last palette card clicked in studio (disambiguates presets with identical colors). */
let studioPaletteId = null;

/** @type {(() => void)|null} */
let studioVisionStripReschedule = null;

/** @type {(() => void)|null} */
let studioStatePersistHandler = null;

/** @type {(() => void)|null} */
let studioStatePanelRefresh = null;

/** @type {CarouselDeck|null} */
let studioDeckRef = null;

/** Right-float strip: render at 80% of 2× center preview; panel fit prioritizes readable height. */
const STRIP_PREVIEW_SCALE = 1.6;
const STRIP_PANEL_WIDTH_PX = PANORAMA_SLIDE_WIDTH_PX;
const STRIP_PANEL_MAX_HEIGHT_PX = 420;
/** Center grid variant editors vs deck previewMaxPx (strip/export keep deck sizing). */
const STUDIO_PREVIEW_SCALE = 0.72;

/** @param {import('./theme.js').DeckTheme} theme */
function studioPreviewMaxPx(theme) {
  return Math.max(180, Math.round(theme.previewMaxPx * STUDIO_PREVIEW_SCALE));
}

/** @param {import('./theme.js').DeckTheme} theme */
function stripDisplayWidthPx(theme) {
  return Math.round(theme.previewMaxPx * STRIP_PREVIEW_SCALE);
}

/** Render width (renderer min/max clamp); may exceed display width. */
/** @param {import('./theme.js').DeckTheme} theme */
function stripRenderWidthPx(theme) {
  const desired = stripDisplayWidthPx(theme);
  return Math.max(SLIDE_SIZE_MIN, Math.min(SLIDE_SIZE_MAX, desired));
}

/** Cap from 80% of 2× target (before panel fit). */
/** @param {import('./theme.js').DeckTheme} theme */
function stripDisplayScale(theme) {
  const render = stripRenderWidthPx(theme);
  return render > 0 ? stripDisplayWidthPx(theme) / render : 1;
}

/**
 * Scale strip canvas for the right float: readable height, horizontal scroll for width.
 * @param {HTMLCanvasElement} strip
 * @param {HTMLElement} mount
 * @param {import('./theme.js').DeckTheme} theme
 */
function resolveStripPanelDisplayScale(strip, mount, theme) {
  const intended = stripDisplayScale(theme);
  const heightScale = STRIP_PANEL_MAX_HEIGHT_PX / strip.height;
  void mount;
  return Math.min(intended, heightScale);
}

function stripPreviewRenderWidthPx(theme) {
  return Math.max(180, stripDisplayWidthPx(theme));
}

/** @param {import('./theme.js').DeckTheme} theme */
function stripPreviewBuildOptions(theme) {
  return {
    renderWidth: stripPreviewRenderWidthPx(theme),
    resolveDisplayScale: (/** @type {HTMLCanvasElement} */ strip, /** @type {HTMLElement} */ mount) =>
      resolveStripPanelDisplayScale(strip, mount, theme),
  };
}

/** @returns {boolean} */
function readShowLineBoxesPreference() {
  return studioShowLineBoxes;
}

/** @param {boolean} enabled */
function writeShowLineBoxesPreference(enabled) {
  studioShowLineBoxes = enabled;
  studioStatePersistHandler?.();
}

let studioShowLineBoxes = false;

/**
 * @typedef {Object} CarouselDeck
 * @property {number} [version]
 * @property {string} [title]
 * @property {string} [slug]
 * @property {Partial<import('./theme.js').DeckTheme>} [deck]
 * @property {Array<{number:number, role?:string, variants:import('./renderer.js').SlideVariant[]}>} slides
 */

/**
 * @param {CarouselDeck} deck
 * @param {{ number: number, role?: string }} slide
 * @param {import('./renderer.js').RenderOptions} renderContext
 * @returns {import('./renderer.js').RenderOptions}
 */
function slideRenderOptions(deck, slide, renderContext) {
  const role = (slide.role || '').trim().toLowerCase();
  return {
    ...renderContext,
    slideRole: role,
    backgroundPanoramaContext: isCarouselCtaRole(role)
      ? undefined
      : buildBackgroundPanoramaContext(deck, slide.number) ?? undefined,
    motifStripContext: isCarouselCtaRole(role) || !motifStripEnabledInDeck(deck)
      ? undefined
      : buildMotifStripContext(deck, slide.number) ?? undefined,
  };
}

/**
 * @typedef {Object} PreviewSlot
 * @property {{number:number}} slide
 * @property {import('./renderer.js').SlideVariant} variant
 * @property {HTMLCanvasElement} preview
 * @property {HTMLElement} previewFrame
 * @property {string} slug
 * @property {() => Promise<void>} refresh
 */

/**
 * @param {{ deckUrl?: string, deck?: CarouselDeck, root?: HTMLElement }} options
 */
export async function mountCarouselStudio(options = {}) {
  const root = options.root || document.getElementById('carousel-app');
  if (!root) throw new Error('Missing carousel root element');

  if ('scrollRestoration' in history) {
    history.scrollRestoration = 'manual';
  }

  root.classList.add('carousel-app');
  root.innerHTML = '<p class="carousel-status">Loading deck…</p>';

  try {
    /** @type {CarouselDeck} */
    let deck;
    if (options.deck) {
      deck = options.deck;
    } else if (options.deckUrl) {
      const cacheBust = options.cacheBust ?? Date.now();
      const separator = options.deckUrl.includes('?') ? '&' : '?';
      const deckRequestUrl = `${options.deckUrl}${separator}v=${cacheBust}`;
      const response = await fetch(deckRequestUrl, { cache: 'no-store' });
      if (!response.ok) {
        throw new Error(`Could not load deck (${response.status})`);
      }
      deck = await response.json();
    } else {
      throw new Error('Provide deck or deckUrl');
    }

    if (!deck.deck) deck.deck = {};
    studioDeckRef = deck;
    const slug = deck.slug || slugify(deck.title || 'carousel');
    const savedScroll = readStudioScrollState(slug);
    let pendingStripScrollLeft = savedScroll?.stripLeft;
    let studioScrollBound = false;
    const fileDeckSnapshot = captureFileDeckSnapshot(deck.deck);
    const fileStripSnapshot = captureFileStripSnapshot(deck);
    const theme = mergeTheme(deck.deck);
    const restoredTheme = applyStudioThemeState(slug, deck, theme);
    applyDeckThemeMargins(theme, deck);
    studioPaletteId = restoredTheme.paletteId ?? matchPaletteId(theme);
    if (restoredTheme.showLineBoxes != null) {
      studioShowLineBoxes = restoredTheme.showLineBoxes;
    }
    await loadDeckFonts(theme);

    root.style.setProperty('--carousel-preview-max', `${studioPreviewMaxPx(theme)}px`);
    root.style.setProperty('--carousel-strip-panel-width', `${STRIP_PANEL_WIDTH_PX}px`);
    root.style.setProperty('--carousel-strip-panel-max-height', `${STRIP_PANEL_MAX_HEIGHT_PX}px`);

    let persistStudioTimer = 0;
    /** @type {(() => void)|null} */
    let refreshStatePanel = null;
    studioStatePersistHandler = () => {
      refreshStatePanel?.();
      window.clearTimeout(persistStudioTimer);
      persistStudioTimer = window.setTimeout(() => {
        writeStudioThemeState(slug, deck, theme, studioPaletteId, studioShowLineBoxes);
      }, 200);
    };
    const bundleBaseUrl = options.bundleBaseUrl ?? resolveBundleBaseUrl(deck, options.deckUrl);
    const assetBaseUrl = options.assetBaseUrl ?? bundleBaseUrl;
    const deckCta = deck.deck?.cta && typeof deck.deck.cta === 'object' ? deck.deck.cta : undefined;
    /** @type {import('./renderer.js').RenderOptions} */
    const renderContext = {
      assetBaseUrl,
      bundleBaseUrl,
      deckCta,
      deck,
    };
    if (deckCta) {
      const ctaTheme = mergeTheme(deck.deck || {});
      try {
        await preloadPostCtaAssets(
          deckCta,
          assetBaseUrl,
          (token) => resolveColor(ctaTheme, token),
          bundleBaseUrl,
        );
      } catch (error) {
        console.warn('[carousel] CTA preload failed:', error);
      }
    }
    /** @param {import('./theme.js').DeckTheme} liveTheme */
    const renderOverrides = (liveTheme) => buildThemeRenderOverrides(deck.deck, liveTheme);
    /** @type {PreviewSlot[]} */
    const previewSlots = [];
    /** @type {Map<number, HTMLInputElement>} */
    const inclusionInputs = new Map();
    /** @type {Map<number, HTMLButtonElement[]>} */
    const stripVariantSegments = new Map();
    /** @type {Set<number>} */
    let includedSlideNumbers = readIncludedSlideNumbers(slug, deck);
    /** @type {Map<number, number>} */
    let stripVariantIndices = readStripVariantIndices(slug, deck);

    /** @returns {{ includedSlideNumbers: Set<number>, variantIndexFor: (slideNumber: number) => number }} */
    const getStudioStripOptions = () => ({
      includedSlideNumbers,
      variantIndexFor: (slideNumber) => stripVariantIndexFor(slug, deck, slideNumber, stripVariantIndices),
    });

    const syncStripVariantUi = () => {
      for (const slide of deck.slides ?? []) {
        const activeIndex = stripVariantIndexFor(slug, deck, slide.number, stripVariantIndices);
        const segments = stripVariantSegments.get(slide.number) ?? [];
        for (let variantIndex = 0; variantIndex < segments.length; variantIndex += 1) {
          const button = segments[variantIndex];
          const active = variantIndex === activeIndex;
          button.classList.toggle('carousel-segment--active', active);
          button.setAttribute('aria-pressed', String(active));
        }
        for (const card of document.querySelectorAll(
          `.carousel-variant[data-slide-number="${slide.number}"]`,
        )) {
          if (!(card instanceof HTMLElement)) continue;
          const variantIndex = Number(card.dataset.variantIndex);
          card.classList.toggle('carousel-variant--in-strip', variantIndex === activeIndex);
        }
      }
    };

    /** @type {ReturnType<typeof installVisionStripController>|null} */
    let visionStripApi = null;
    let visionStripRebuildTimer = 0;

    const syncInclusionInputs = () => {
      for (const [slideNumber, input] of inclusionInputs) {
        input.checked = includedSlideNumbers.has(slideNumber);
        input.indeterminate = false;
      }
    };

    /** @param {number|null} [focusSlideNumber] Scroll strip to this slide after rebuild. */
    const scheduleVisionStripRebuild = (focusSlideNumber = null) => {
      window.clearTimeout(visionStripRebuildTimer);
      visionStripRebuildTimer = window.setTimeout(async () => {
        if (!visionStripApi) return;
        const mount = document.querySelector('.carousel-vision-strip-mount:not(.carousel-vision-strip-mount--hidden)');
        let savedScrollLeft = 0;
        if (mount instanceof HTMLElement) {
          savedScrollLeft = mount.scrollLeft;
          if (focusSlideNumber == null && pendingStripScrollLeft != null) {
            savedScrollLeft = pendingStripScrollLeft;
          }
          showVisionStripPlaceholder(mount);
        }
        try {
          const buildOptions = stripPreviewBuildOptions(theme);
          const strip = await visionStripApi.build(buildOptions.renderWidth, buildOptions);
          if (mount instanceof HTMLElement && strip instanceof HTMLCanvasElement) {
            restoreVisionStripScroll(mount, strip, {
              focusSlideNumber,
              scrollLeft: focusSlideNumber == null ? savedScrollLeft : undefined,
            });
            if (focusSlideNumber == null) {
              pendingStripScrollLeft = undefined;
            }
          }
        } catch (error) {
          console.error('[carousel] vision strip rebuild failed:', error);
          if (mount instanceof HTMLElement) {
            showVisionStripPlaceholder(
              mount,
              error instanceof Error ? error.message : String(error),
            );
          }
        } finally {
          if (!studioScrollBound) {
            applyStudioScrollState(root, savedScroll, {
              window: false,
              leftFloat: false,
              rightFloat: true,
              strip: false,
            });
            bindStudioScrollPersistence(slug, root);
            studioScrollBound = true;
          }
        }
      }, 120);
    };

    /** @param {number|null} [focusSlideNumber] */
    const applySlideInclusionChange = async (focusSlideNumber = null) => {
      writeIncludedSlideNumbers(slug, includedSlideNumbers);
      syncInclusionInputs();
      studioStatePanelRefresh?.();
      await refreshAllPreviews(previewSlots);
      scheduleVisionStripRebuild(focusSlideNumber);
    };

    /** @param {number} slideNumber */
    const applyStripVariantChange = (slideNumber) => {
      writeStripVariantIndices(slug, stripVariantIndices);
      syncStripVariantUi();
      studioStatePanelRefresh?.();
      scheduleVisionStripRebuild(slideNumber);
    };

    root.innerHTML = '';

    root.appendChild(renderLeftFloatPanel(theme, deck, previewSlots));

    const rightFloat = renderRightFloatPanel();
    root.appendChild(rightFloat);

    const visionStripMount = rightFloat.querySelector('.carousel-vision-strip-mount');
    if (!(visionStripMount instanceof HTMLElement)) throw new Error('Missing vision strip mount');
    showVisionStripPlaceholder(visionStripMount);

    visionStripApi = installVisionStripController(
      deck,
      theme,
      renderOverrides,
      renderContext,
      visionStripMount,
      getStudioStripOptions,
    );
    wireVisionStripPanel(rightFloat, visionStripApi, theme);

    const statePanel = renderStudioStatePanel({
      fileDeckSnapshot,
      fileStripSnapshot,
      getDeck: () => studioDeckRef ?? deck,
      getTheme: () => theme,
      getPaletteId: () => studioPaletteId,
      getShowLineBoxes: () => studioShowLineBoxes,
      getIncludedSlideNumbers: () => includedSlideNumbers,
      getStripVariantIndices: () => stripVariantIndices,
      onSave: async () => {
        writeStudioThemeState(slug, deck, theme, studioPaletteId, studioShowLineBoxes);
        if (deck.deck && theme.backgroundWave) {
          deck.deck.backgroundWave = compactBackgroundWaveForExport(theme.backgroundWave);
          deck.deck.palette = deckPaletteFromTheme(theme);
        }
        const mode = await saveCarouselDeckToSource(deck, { deckUrl: options.deckUrl });
        const next = captureFileDeckSnapshot(deck.deck);
        for (const key of Object.keys(fileDeckSnapshot)) {
          delete fileDeckSnapshot[key];
        }
        Object.assign(fileDeckSnapshot, next);
        refreshStatePanel?.();
        if (mode === 'download') {
          window.alert('Saved a copy to Downloads. Replace the bundle carousel.json with that file, or restart with make serve so Save can write the source file.');
        }
      },
    });
    refreshStatePanel = statePanel.refresh;
    studioStatePanelRefresh = statePanel.refresh;
    rightFloat.appendChild(statePanel.element);

    root.appendChild(renderSlides(
      deck,
      theme,
      slug,
      previewSlots,
      renderOverrides,
      renderContext,
      inclusionInputs,
      includedSlideNumbers,
      stripVariantIndices,
      stripVariantSegments,
      (slideNumber, included) => {
        if (included) {
          includedSlideNumbers.add(slideNumber);
        } else {
          includedSlideNumbers.delete(slideNumber);
        }
        applySlideInclusionChange(slideNumber);
      },
      (slideNumber, variantIndex) => {
        stripVariantIndices.set(slideNumber, variantIndex);
        applyStripVariantChange(slideNumber);
      },
    ));

    await refreshAllPreviews(previewSlots);
    attachPreviewResizeObservers(previewSlots);
    syncStripVariantUi();
    studioVisionStripReschedule = scheduleVisionStripRebuild;
    applyStudioScrollState(root, savedScroll, {
      window: true,
      leftFloat: true,
      rightFloat: false,
      strip: false,
    });
    requestAnimationFrame(() => {
      applyStudioScrollState(root, savedScroll, {
        window: true,
        leftFloat: true,
        rightFloat: false,
        strip: false,
      });
    });
    scheduleVisionStripRebuild();
  } catch (error) {
    root.innerHTML = '';
    const box = document.createElement('div');
    box.className = 'carousel-error';
    box.innerHTML = `<strong>Carousel preview failed</strong><p>${escapeHtml(error instanceof Error ? error.message : String(error))}</p><p>Run <code>make serve</code> and open this page through Hugo (not file://).</p>`;
    root.appendChild(box);
  }
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {CarouselDeck} deck
 * @param {PreviewSlot[]} previewSlots
 */
function renderLeftFloatPanel(theme, deck, previewSlots) {
  const panel = document.createElement('aside');
  panel.className = 'carousel-float carousel-float--left';
  panel.setAttribute('aria-label', 'Palette, wave, and layout debug');

  const scrollBody = document.createElement('div');
  scrollBody.className = 'carousel-float-body';

  const themeSection = document.createElement('div');
  themeSection.className = 'carousel-float-theme';
  themeSection.appendChild(renderPalettePanel(theme, previewSlots));
  themeSection.appendChild(renderGradientPanel(theme, previewSlots));
  scrollBody.appendChild(themeSection);

  const debugSection = document.createElement('div');
  debugSection.className = 'carousel-float-debug';
  if (deckHasMotifStrip(deck)) {
    debugSection.appendChild(renderMotifStripPanel(deck, previewSlots));
  }
  if (deckHasPostCta(deck)) {
    debugSection.appendChild(renderCtaSettingsPanel(deck, previewSlots));
  }
  debugSection.appendChild(renderMarginsPanel(theme, deck, previewSlots));
  debugSection.appendChild(renderLineHeightsPanel(theme, deck, previewSlots));
  debugSection.appendChild(renderLineBoxesPanel(previewSlots));
  scrollBody.appendChild(debugSection);

  panel.appendChild(scrollBody);

  return panel;
}

/** Right-float strip preview panel. */
function renderRightFloatPanel() {
  const panel = document.createElement('aside');
  panel.id = VISION_STRIP_PANEL_ID;
  panel.className = 'carousel-float carousel-float--right';
  panel.setAttribute('aria-label', 'Carousel strip preview');

  const stripSection = document.createElement('section');
  stripSection.className = 'carousel-strip-preview';

  const mount = document.createElement('div');
  mount.className = 'carousel-vision-strip-mount';

  const actions = document.createElement('div');
  actions.className = 'carousel-vision-strip-actions';

  const downloadActions = document.createElement('div');
  downloadActions.className = 'carousel-vision-strip-downloads';
  downloadActions.append(
    createLabeledDownloadButton({
      id: 'carousel-vision-strip-download',
      label: 'Download panorama',
      title: 'Export the strip as one panoramic WebP image at full slide resolution',
      text: 'Panorama',
    }),
    createLabeledDownloadButton({
      id: 'carousel-vision-strip-download-all',
      label: 'Download all slides separately',
      title: 'Export each slide in the strip as its own WebP file at full resolution',
      text: 'Slides',
    }),
    createLabeledDownloadButton({
      id: 'carousel-vision-strip-download-pdf',
      label: 'Download LinkedIn PDF',
      title: 'Export In-strip slides as one full-bleed LinkedIn PDF',
      text: 'PDF',
    }),
  );

  actions.append(downloadActions);
  stripSection.append(mount, actions);
  panel.appendChild(stripSection);

  return panel;
}

/**
 * @param {{ id: string, label: string, title?: string, text: string }} options
 * @returns {HTMLButtonElement}
 */
function createLabeledDownloadButton(options) {
  const btn = createDownloadButton({
    id: options.id,
    className: 'carousel-button carousel-button--compact carousel-download-btn--labeled',
    label: options.label,
    title: options.title ?? options.label,
  });
  const label = document.createElement('span');
  label.className = 'carousel-download-btn-label';
  label.textContent = options.text;
  btn.appendChild(label);
  return btn;
}

/** @param {PreviewSlot[]} previewSlots */
function renderLineBoxesPanel(previewSlots) {
  const panel = document.createElement('div');
  panel.className = 'carousel-debug-lineboxes';
  panel.appendChild(renderLineBoxesToggle(previewSlots));
  return panel;
}

/**
 * @param {CarouselDeck} deck
 * @param {PreviewSlot[]} previewSlots
 */
function renderMotifStripPanel(deck, previewSlots) {
  const panel = document.createElement('div');
  panel.className = 'carousel-debug-motif';

  const label = document.createElement('label');
  label.className = 'carousel-debug-toggle';

  const input = document.createElement('input');
  input.type = 'checkbox';
  input.checked = motifStripEnabledInDeck(deck);
  input.setAttribute('aria-label', 'Enable motif strip on slides');
  input.addEventListener('change', () => {
    setMotifStripEnabled(deck, input.checked);
    refreshAllPreviews(previewSlots);
    studioStatePersistHandler?.();
    studioStatePanelRefresh?.();
  });

  const spec = /** @type {Record<string, unknown>|undefined} */ (
    deck.deck?.motifStrip && typeof deck.deck.motifStrip === 'object' && !Array.isArray(deck.deck.motifStrip)
      ? deck.deck.motifStrip
      : undefined
  );
  const src = typeof spec?.src === 'string' ? spec.src.trim() : '';

  label.appendChild(input);
  label.append(document.createTextNode(' Motif strip'));

  if (src) {
    const hint = document.createElement('span');
    hint.className = 'carousel-debug-toggle-hint';
    hint.textContent = `Writes motifStrip.enabled in carousel.json (${src})`;
    label.appendChild(hint);
  }

  panel.appendChild(label);

  if (src) {
    colorControlRefs.motifOffsetInputs = [];

    const scaleField = document.createElement('div');
    scaleField.className = 'carousel-lineheight-field';
    const scaleLabel = document.createElement('span');
    scaleLabel.className = 'carousel-lineheight-label';
    scaleLabel.textContent = 'Size';
    const scaleWrap = document.createElement('span');
    scaleWrap.className = 'carousel-lineheight-input-wrap carousel-motif-scale-wrap';
    const minusBtn = document.createElement('button');
    minusBtn.type = 'button';
    minusBtn.className = 'carousel-motif-scale-btn';
    minusBtn.textContent = '−';
    minusBtn.setAttribute('aria-label', 'Decrease motif size by 1 percent');
    const scaleInput = document.createElement('input');
    scaleInput.type = 'text';
    scaleInput.inputMode = 'decimal';
    scaleInput.autocomplete = 'off';
    scaleInput.spellcheck = false;
    scaleInput.className = 'carousel-lineheight-input carousel-motif-scale-input';
    scaleInput.value = formatMotifScaleDelta(motifScaleDelta(deck));
    scaleInput.setAttribute('aria-label', 'Motif size change from native, as a signed percent');
    const scaleSuffix = document.createElement('span');
    scaleSuffix.className = 'carousel-lineheight-suffix';
    scaleSuffix.textContent = '%';
    scaleSuffix.setAttribute('aria-hidden', 'true');
    const plusBtn = document.createElement('button');
    plusBtn.type = 'button';
    plusBtn.className = 'carousel-motif-scale-btn';
    plusBtn.textContent = '+';
    plusBtn.setAttribute('aria-label', 'Increase motif size by 1 percent');

    const applyScaleDelta = (raw) => {
      const parsed = parseMotifScaleDelta(raw);
      if (parsed == null) return;
      setMotifStripScaleDelta(deck, parsed);
      scaleInput.value = formatMotifScaleDelta(motifScaleDelta(deck));
      refreshAllPreviews(previewSlots);
      studioStatePersistHandler?.();
      studioStatePanelRefresh?.();
    };

    const nudgeScale = (step) => {
      applyScaleDelta(motifScaleDelta(deck) + step);
    };

    scaleInput.addEventListener('change', () => {
      const parsed = parseMotifScaleDelta(scaleInput.value);
      applyScaleDelta(parsed == null ? motifScaleDelta(deck) : parsed);
    });
    scaleInput.addEventListener('keydown', (event) => {
      if (event.key === 'Enter') {
        event.preventDefault();
        scaleInput.blur();
      }
    });
    minusBtn.addEventListener('click', () => nudgeScale(-1));
    plusBtn.addEventListener('click', () => nudgeScale(1));
    scaleInput._refreshMotifOffset = () => {
      scaleInput.value = formatMotifScaleDelta(motifScaleDelta(deck));
    };
    scaleWrap.append(minusBtn, scaleInput, scaleSuffix, plusBtn);
    scaleField.append(scaleLabel, scaleWrap);
    panel.appendChild(scaleField);
    colorControlRefs.motifOffsetInputs.push(scaleInput);

    const fields = document.createElement('div');
    fields.className = 'carousel-lineheight-fields';

    /** @type {readonly { key: 'offsetX' | 'offsetY', label: string }[]} */
    const entries = [
      { key: 'offsetX', label: 'Horizontal' },
      { key: 'offsetY', label: 'Vertical' },
    ];

    for (const entry of entries) {
      const field = document.createElement('div');
      field.className = 'carousel-lineheight-field';

      const fieldLabel = document.createElement('span');
      fieldLabel.className = 'carousel-lineheight-label';
      fieldLabel.textContent = entry.label;

      const input = document.createElement('input');
      input.type = 'number';
      input.className = 'carousel-lineheight-input';
      const offsetMax = motifOffsetMax(entry.key);
      input.min = String(-offsetMax);
      input.max = String(offsetMax);
      input.step = '1';
      input.dataset.motifOffsetKey = entry.key;
      input.value = String(motifOffsetValue(deck, entry.key));
      input.setAttribute('aria-label', `${entry.label} motif offset in design pixels at 1080`);

      input.addEventListener('input', () => {
        const parsed = Number(input.value);
        if (!Number.isFinite(parsed)) return;
        setMotifStripOffset(deck, entry.key, parsed);
        refreshAllPreviews(previewSlots);
        studioStatePersistHandler?.();
        studioStatePanelRefresh?.();
      });

      input.addEventListener('change', () => {
        const parsed = Number(input.value);
        const next = Number.isFinite(parsed) ? parsed : motifOffsetValue(deck, entry.key);
        setMotifStripOffset(deck, entry.key, next);
        input.value = String(motifOffsetValue(deck, entry.key));
        refreshAllPreviews(previewSlots);
        studioStatePersistHandler?.();
        studioStatePanelRefresh?.();
      });

      input._refreshMotifOffset = () => {
        input.value = String(motifOffsetValue(deck, entry.key));
      };

      field.append(fieldLabel, input);
      fields.appendChild(field);
      colorControlRefs.motifOffsetInputs.push(input);
    }

    panel.appendChild(fields);

    const offsetHint = document.createElement('p');
    offsetHint.className = 'carousel-debug-lineheights-hint';
    offsetHint.textContent = 'Size scales width and height together from the bottom-left (0 is unchanged, -10 is 10% smaller, +10 is 10% larger). The strip stays one continuous line across slides. Writes bandWidth, offsetX, offsetY.';
    panel.appendChild(offsetHint);
  }

  return panel;
}

/** @param {PreviewSlot[]} previewSlots */
function renderLineBoxesToggle(previewSlots) {
  const label = document.createElement('label');
  label.className = 'carousel-debug-toggle';

  const input = document.createElement('input');
  input.type = 'checkbox';
  input.checked = studioShowLineBoxes;
  input.addEventListener('change', () => {
    studioShowLineBoxes = input.checked;
    writeShowLineBoxesPreference(studioShowLineBoxes);
    refreshAllPreviews(previewSlots);
  });

  const hint = document.createElement('span');
  hint.className = 'carousel-debug-toggle-hint';
  hint.textContent = 'Blue = line slot; orange = ascender + descender; red = baseline + x-height top';

  label.appendChild(input);
  label.append(document.createTextNode(' Line boxes'));
  label.appendChild(hint);

  return label;
}

/**
 * @param {ReturnType<typeof installVisionStripController>} visionStripApi
 * @param {import('./theme.js').DeckTheme} theme
 */
function wireVisionStripPanel(panel, visionStripApi, theme) {
  wireStripDownloadButton(panel, '#carousel-vision-strip-download', async () => {
    await visionStripApi.download({ slideWidth: theme.size });
  });
  wireStripDownloadButton(panel, '#carousel-vision-strip-download-all', async () => {
    await visionStripApi.downloadAll();
  });
  wireStripDownloadButton(panel, '#carousel-vision-strip-download-pdf', async () => {
    await visionStripApi.downloadPdf();
  });
}

/**
 * @param {HTMLElement} panel
 * @param {string} selector
 * @param {() => Promise<void>} handler
 */
function wireStripDownloadButton(panel, selector, handler) {
  const downloadBtn = panel.querySelector(selector);
  if (!(downloadBtn instanceof HTMLButtonElement)) return;

  downloadBtn.addEventListener('click', async () => {
    downloadBtn.disabled = true;
    downloadBtn.setAttribute('aria-busy', 'true');
    downloadBtn.classList.add('carousel-icon-button--busy');
    try {
      await handler();
    } catch (error) {
      console.error('[carousel] strip download failed:', error);
      window.alert(error instanceof Error ? error.message : String(error));
    }
    downloadBtn.disabled = false;
    downloadBtn.removeAttribute('aria-busy');
    downloadBtn.classList.remove('carousel-icon-button--busy');
  });
}

/**
 * @param {CarouselDeck} deck
 * @param {import('./theme.js').DeckTheme} theme
 * @param {string} slug
 * @param {PreviewSlot[]} previewSlots
 * @param {(live: import('./theme.js').DeckTheme) => Partial<import('./theme.js').DeckTheme> & Record<string, unknown>} renderOverrides
 * @param {import('./renderer.js').RenderOptions} renderContext
 */
function renderSlides(
  deck,
  theme,
  slug,
  previewSlots,
  renderOverrides,
  renderContext,
  inclusionInputs,
  includedSlideNumbers,
  stripVariantIndices,
  stripVariantSegments,
  onInclusionChange,
  onStripVariantChange,
) {
  const list = document.createElement('div');
  list.className = 'carousel-slides';

  for (const slide of deck.slides) {
    const card = document.createElement('section');
    card.className = 'carousel-slide';

    const head = document.createElement('div');
    head.className = 'carousel-slide-head';

    const title = document.createElement('h2');
    title.textContent = `Slide ${slide.number}`;

    const role = document.createElement('span');
    role.className = 'carousel-slide-role';
    role.textContent = slide.role || 'slide';

    const inclusionLabel = document.createElement('label');
    inclusionLabel.className = 'carousel-slide-inclusion';

    const inclusionInput = document.createElement('input');
    inclusionInput.type = 'checkbox';
    inclusionInput.checked = includedSlideNumbers.has(slide.number);
    inclusionInput.title = 'Include this slide in the strip preview';
    inclusionInput.addEventListener('change', () => {
      onInclusionChange(slide.number, inclusionInput.checked);
    });

    inclusionLabel.append(inclusionInput, document.createTextNode(' In strip'));
    inclusionInputs.set(slide.number, inclusionInput);

    const stripControls = document.createElement('div');
    stripControls.className = 'carousel-slide-strip-controls';
    stripControls.appendChild(inclusionLabel);

    if (slide.variants.length > 1) {
      const pick = document.createElement('div');
      pick.className = 'carousel-slide-strip-pick';

      const pickLabel = document.createElement('span');
      pickLabel.className = 'carousel-slide-strip-pick-label';
      pickLabel.textContent = 'Variant';

      const group = document.createElement('div');
      group.className = 'carousel-segmented carousel-segmented--compact';
      group.setAttribute('role', 'group');
      group.setAttribute('aria-label', `Slide ${slide.number} strip variant`);

      /** @type {HTMLButtonElement[]} */
      const segments = [];
      const activeIndex = stripVariantIndices.get(slide.number) ?? 0;

      for (let variantIndex = 0; variantIndex < slide.variants.length; variantIndex += 1) {
        const label = variantIdFromIndex(variantIndex);
        const btn = document.createElement('button');
        btn.type = 'button';
        btn.className = 'carousel-segment carousel-segment--compact';
        btn.textContent = label;
        btn.title = `Use variant ${label} in the strip preview`;
        btn.setAttribute('aria-pressed', String(activeIndex === variantIndex));
        if (activeIndex === variantIndex) {
          btn.classList.add('carousel-segment--active');
        }
        btn.addEventListener('click', () => {
          onStripVariantChange(slide.number, variantIndex);
        });
        segments.push(btn);
        group.appendChild(btn);
      }

      stripVariantSegments.set(slide.number, segments);
      pick.append(pickLabel, group);
      stripControls.appendChild(pick);
    }

    head.append(title, role, stripControls);

    card.appendChild(head);

    const variants = document.createElement('div');
    variants.className = 'carousel-variants';

    for (let variantIndex = 0; variantIndex < slide.variants.length; variantIndex += 1) {
      variants.appendChild(
        renderVariantCard(
          deck,
          slide,
          slide.variants[variantIndex],
          variantIndex,
          theme,
          slug,
          previewSlots,
          renderOverrides,
          renderContext,
          stripVariantIndices.get(slide.number) === variantIndex,
        ),
      );
    }

    card.appendChild(variants);
    list.appendChild(card);
  }

  return list;
}

/**
 * @param {CarouselDeck} deck
 * @param {{number:number, role?:string}} slide
 * @param {import('./renderer.js').SlideVariant} variant
 * @param {number} variantIndex 0-based index in `slide.variants` (export label `a`, `b`, …)
 * @param {import('./theme.js').DeckTheme} theme
 * @param {string} slug
 * @param {PreviewSlot[]} previewSlots
 * @param {(live: import('./theme.js').DeckTheme) => Partial<import('./theme.js').DeckTheme> & Record<string, unknown>} renderOverrides
 */
function renderVariantCard(
  deck,
  slide,
  variant,
  variantIndex,
  theme,
  slug,
  previewSlots,
  renderOverrides,
  renderContext,
  isStripVariant,
) {
  const variantLabel = variantIdFromIndex(variantIndex);
  const card = document.createElement('article');
  card.className = 'carousel-variant';
  card.dataset.slideNumber = String(slide.number);
  card.dataset.variantIndex = String(variantIndex);
  if (isStripVariant) {
    card.classList.add('carousel-variant--in-strip');
  }

  const previewFrame = document.createElement('div');
  previewFrame.className = 'carousel-preview-frame';

  const preview = document.createElement('canvas');
  preview.className = 'carousel-preview';
  previewFrame.appendChild(preview);

  /** @type {PreviewSlot} */
  const slot = {
    slide,
    variant,
    preview,
    previewFrame,
    slug,
    refresh: async () => {},
  };

  slot.refresh = async () => {
    await new Promise((resolve) => requestAnimationFrame(resolve));
    const overrides = renderOverrides(theme);
    const aspect = parseAspectRatio(overrides.aspectRatio);
    const { width: pixelWidth, height: pixelHeight } = previewPixelDimensions(
      previewFrame,
      aspect,
    );
    preview.width = pixelWidth;
    preview.height = pixelHeight;
    preview.style.width = '100%';
    preview.style.height = '100%';
    await renderSlideToCanvas(variant, overrides, {
      grain: false,
      studioPreview: true,
      targetCanvas: preview,
      outputSize: pixelWidth,
      showLineBoxes: studioShowLineBoxes,
      ...slideRenderOptions(deck, slide, renderContext),
    });
  };

  const toolbar = document.createElement('div');
  toolbar.className = 'carousel-variant-toolbar';

  const alignmentCopyBtn = createCopyButton({
    className: 'carousel-control-copy',
    label: 'Copy placement',
    title: 'Copy verticalAlign and alignment JSON for this variant',
  });
  alignmentCopyBtn.addEventListener('click', () => {
    copyText(alignmentCopyBtn, placementJsonSnippet(variant), 'Copy placement');
  });

  const exportBtn = createDownloadButton({
    className: 'carousel-control-copy',
    label: `Export variant ${variantLabel} WebP`,
  });
  exportBtn.addEventListener('click', async () => {
    exportBtn.disabled = true;
    const canvas = await renderSlideToCanvas(variant, renderOverrides(theme), {
      supersample: 2,
      ...slideRenderOptions(deck, slide, renderContext),
    });
    downloadCanvasWebp(canvas, exportFilename(slug, slide.number, variantIndex));
    exportBtn.disabled = false;
  });

  toolbar.appendChild(alignmentCopyBtn);
  toolbar.appendChild(exportBtn);

  const body = document.createElement('div');
  body.className = 'carousel-variant-body';
  // post_cta uses a fixed top-down funnel layout; verticalAlign has no effect.
  if (variant.archetype !== 'post_cta') {
    body.appendChild(renderVariantVerticalAlignRail(slide, variant, variantLabel, slot));
  }
  body.appendChild(previewFrame);
  body.appendChild(renderVariantSideControls(slide, variant, variantLabel, slot));

  card.appendChild(toolbar);
  card.appendChild(body);

  previewSlots.push(slot);

  return card;
}

/**
 * @param {{ number: number }} slide
 * @param {import('./renderer.js').SlideVariant} variant
 * @param {string} variantLabel
 * @param {PreviewSlot} slot
 */
function renderVariantVerticalAlignRail(slide, variant, variantLabel, slot) {
  const rail = document.createElement('aside');
  rail.className = 'carousel-variant-side-controls carousel-variant-align-rail';

  const zone = document.createElement('div');
  zone.className = 'carousel-side-slot carousel-side-slot--body';
  zone.appendChild(renderSideControl({
    slide,
    variant,
    variantLabel,
    slot,
    axis: 'vertical',
    label: 'Vert',
    values: [
      { value: 'top', text: 'T', title: 'Top' },
      { value: 'center', text: 'M', title: 'Middle' },
      { value: 'bottom', text: 'B', title: 'Bottom' },
    ],
    current: () => currentVerticalAlign(variant),
    apply: (value) => setVerticalAlign(variant, /** @type {'top'|'center'|'bottom'} */ (value)),
  }));
  rail.appendChild(zone);

  return rail;
}

/**
 * @param {{ number: number }} slide
 * @param {import('./renderer.js').SlideVariant} variant
 * @param {string} variantLabel
 * @param {PreviewSlot} slot
 */
function renderVariantSideControls(slide, variant, variantLabel, slot) {
  const rail = document.createElement('aside');
  rail.className = 'carousel-variant-side-controls';

  for (const section of sectionsInVariant(variant)) {
    const zone = document.createElement('div');
    zone.className = `carousel-side-slot carousel-side-slot--${section}`;
    zone.appendChild(renderSideControl({
      slide,
      variant,
      variantLabel,
      slot,
      axis: 'horizontal',
      label: sectionControlLabel(section),
      values: [
        { value: 'left', text: 'L', title: 'Left' },
        { value: 'center', text: 'C', title: 'Center' },
        { value: 'right', text: 'R', title: 'Right' },
      ],
      current: () => sectionAlignmentFor(variant, section),
      apply: (value) => setSectionAlignment(variant, section, /** @type {'left'|'center'|'right'} */ (value)),
    }));
    rail.appendChild(zone);
  }

  return rail;
}

/**
 * @param {Object} options
 * @param {{ number: number }} options.slide
 * @param {import('./renderer.js').SlideVariant} options.variant
 * @param {PreviewSlot} options.slot
 * @param {string} [options.label]
 * @param {string} options.ariaLabel
 * @param {Array<{ value: string, text: string, title?: string }>} options.values
 * @param {() => string} options.current
 * @param {(value: string) => void} options.apply
 * @param {() => string} [options.jsonSnippet]
 * @param {string} [options.copyTitle]
 */
function renderControlRow(options) {
  const row = document.createElement('div');
  row.className = 'carousel-variant-control-row';

  if (options.label) {
    const label = document.createElement('span');
    label.className = 'carousel-control-label';
    label.textContent = options.label;
    row.appendChild(label);
  }

  const group = document.createElement('div');
  group.className = 'carousel-segmented';
  group.setAttribute('role', 'group');
  group.setAttribute('aria-label', options.ariaLabel);

  /** @type {HTMLButtonElement[]} */
  const segments = [];
  const activeValue = options.current();

  for (const entry of options.values) {
    const btn = document.createElement('button');
    btn.type = 'button';
    btn.className = 'carousel-segment';
    btn.textContent = entry.text;
    btn.title = entry.title || entry.text;
    btn.dataset.value = entry.value;
    btn.setAttribute('aria-pressed', String(activeValue === entry.value));
    if (activeValue === entry.value) {
      btn.classList.add('carousel-segment--active');
    }
    btn.addEventListener('click', () => {
      options.apply(entry.value);
      for (const segment of segments) {
        const active = segment.dataset.value === entry.value;
        segment.classList.toggle('carousel-segment--active', active);
        segment.setAttribute('aria-pressed', String(active));
      }
      options.slot.refresh();
    });
    segments.push(btn);
    group.appendChild(btn);
  }

  const actions = document.createElement('div');
  actions.className = 'carousel-variant-control-actions';
  actions.appendChild(group);
  if (options.jsonSnippet && options.copyTitle) {
    const copyBtn = createCopyButton({
      className: 'carousel-control-copy',
      label: options.copyTitle,
    });
    copyBtn.addEventListener('click', () => {
      copyText(copyBtn, options.jsonSnippet(), options.copyTitle);
    });
    actions.appendChild(copyBtn);
  }

  row.appendChild(actions);
  return row;
}

/**
 * @param {Object} options
 * @param {{ number: number }} options.slide
 * @param {import('./renderer.js').SlideVariant} options.variant
 * @param {string} options.variantLabel
 * @param {PreviewSlot} options.slot
 * @param {'horizontal' | 'vertical'} options.axis
 * @param {string} options.label
 * @param {Array<{ value: string, text: string, title?: string }>} options.values
 * @param {() => string} options.current
 * @param {(value: string) => void} options.apply
 */
function renderSideControl(options) {
  const block = document.createElement('div');
  block.className = 'carousel-side-control';

  const label = document.createElement('span');
  label.className = 'carousel-side-control-label';
  label.textContent = options.label;

  const group = document.createElement('div');
  group.className = 'carousel-segmented carousel-segmented--side';
  group.setAttribute('role', 'group');
  group.setAttribute(
    'aria-label',
    `Slide ${options.slide.number} variant ${options.variantLabel} ${options.label} ${options.axis}`,
  );

  /** @type {HTMLButtonElement[]} */
  const segments = [];
  const activeValue = options.current();

  for (const entry of options.values) {
    const btn = document.createElement('button');
    btn.type = 'button';
    btn.className = 'carousel-segment carousel-segment--compact';
    btn.textContent = entry.text;
    btn.title = entry.title || entry.text;
    btn.dataset.value = entry.value;
    btn.setAttribute('aria-pressed', String(activeValue === entry.value));
    if (activeValue === entry.value) {
      btn.classList.add('carousel-segment--active');
    }
    btn.addEventListener('click', () => {
      options.apply(entry.value);
      for (const segment of segments) {
        const active = segment.dataset.value === entry.value;
        segment.classList.toggle('carousel-segment--active', active);
        segment.setAttribute('aria-pressed', String(active));
      }
      options.slot.refresh();
    });
    segments.push(btn);
    group.appendChild(btn);
  }

  block.append(label, group);
  return block;
}

/** @param {import('./renderer.js').SlideVariant} variant @returns {import('./theme.js').Section[]} */
function sectionsInVariant(variant) {
  const present = new Set(variant.blocks.map((block) => block.section));
  const sections = /** @type {import('./theme.js').Section[]} */ (
    ['header', 'body', 'footer'].filter((section) => present.has(section))
  );
  // post_cta split layout pins footer copy left in the right column; alignment has no effect.
  if (variant.archetype === 'post_cta') {
    return sections.filter((section) => section !== 'footer');
  }
  return sections;
}

/** @param {import('./theme.js').Section} section */
function sectionControlLabel(section) {
  if (section === 'header') return 'Header';
  if (section === 'footer') return 'Footer';
  return 'Body';
}

/**
 * @param {import('./renderer.js').SlideVariant} variant
 * @param {import('./theme.js').Section} section
 * @param {'left' | 'center' | 'right'} value
 */
function setSectionAlignment(variant, section, value) {
  const normalized = normalizeSectionAlignment(variant);
  normalized[section] = value;
  const compact = compactSectionAlignment({ ...variant, alignment: normalized });
  if (Object.keys(compact).length === 0) {
    delete variant.alignment;
    return;
  }
  variant.alignment = compact;
}

/**
 * Paste-ready variant fields: horizontal `alignment` and optional `verticalAlign`.
 * @param {import('./renderer.js').SlideVariant} variant
 */
function placementJsonSnippet(variant) {
  /** @type {Record<string, unknown>} */
  const payload = {};
  const alignment = compactSectionAlignment(variant);
  if (Object.keys(alignment).length > 0) {
    payload.alignment = alignment;
  }
  if (variant.verticalAlign === 'center' || variant.verticalAlign === 'bottom') {
    payload.verticalAlign = variant.verticalAlign;
  }
  if (Object.keys(payload).length === 0) {
    return '// default placement — omit alignment and verticalAlign';
  }
  return JSON.stringify(payload, null, 2);
}

/** @param {import('./renderer.js').SlideVariant} variant */
function currentVerticalAlign(variant) {
  const v = variant.verticalAlign;
  if (v === 'center' || v === 'bottom') return v;
  return 'top';
}

/**
 * @param {import('./renderer.js').SlideVariant} variant
 * @param {'top' | 'center' | 'bottom'} value
 */
function setVerticalAlign(variant, value) {
  if (value === 'top') {
    delete variant.verticalAlign;
    return;
  }
  variant.verticalAlign = value;
}

/**
 * @param {number} value
 * @param {number} min
 * @param {number} max
 */
function clampLineHeight(value, min, max) {
  return Math.max(min, Math.min(max, value));
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {string} key
 */
function lineHeightValueFor(theme, key) {
  const map = theme.lineHeights || {};
  const raw = map[key];
  return Number.isFinite(raw) ? raw : 1;
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {CarouselDeck} deck
 * @param {string} key
 * @param {number} value
 */
function setDeckLineHeight(theme, deck, key, value) {
  const next = clampLineHeight(value, 0.25, 4);
  if (!theme.lineHeights) {
    theme.lineHeights = {};
  }
  theme.lineHeights[key] = next;
  if (!deck.deck) {
    deck.deck = {};
  }
  if (!deck.deck.lineHeights) {
    deck.deck.lineHeights = {};
  }
  deck.deck.lineHeights[key] = next;
}

/** @param {import('./theme.js').DeckTheme} theme @param {'horizontal' | 'vertical'} axis */
function marginPercentFromTheme(theme, axis) {
  const px = axis === 'horizontal' ? theme.marginHorizontal : theme.marginVertical;
  return Math.round((px / DESIGN_CANVAS_SIZE) * 100);
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {CarouselDeck} deck
 */
function applyDeckThemeMargins(theme, deck) {
  const merged = mergeTheme(deck.deck || {});
  theme.marginHorizontal = merged.marginHorizontal;
  theme.marginVertical = merged.marginVertical;
}

/**
 * Split `deck.margin` shorthand into per-axis % fields before the first axis edit.
 * @param {CarouselDeck} deck
 * @param {import('./theme.js').DeckTheme} theme
 */
function ensurePerAxisMargins(deck, theme) {
  if (!deck.deck?.margin) return;
  if (!deck.deck) deck.deck = {};
  deck.deck.marginHorizontal = `${marginPercentFromTheme(theme, 'horizontal')}%`;
  deck.deck.marginVertical = `${marginPercentFromTheme(theme, 'vertical')}%`;
  delete deck.deck.margin;
}

/**
 * @param {CarouselDeck} deck
 * @param {import('./theme.js').DeckTheme} theme
 * @param {'horizontal' | 'vertical'} axis
 * @param {number} percent
 */
function setDeckMarginPercent(deck, theme, axis, percent) {
  if (!deck.deck) deck.deck = {};
  ensurePerAxisMargins(deck, theme);
  const clamped = Math.round(
    Math.max(MARGIN_PERCENT_MIN, Math.min(MARGIN_PERCENT_MAX, percent)),
  );
  const key = axis === 'horizontal' ? 'marginHorizontal' : 'marginVertical';
  deck.deck[key] = `${clamped}%`;
  applyDeckThemeMargins(theme, deck);
}

/**
 * Section title row with optional copy-to-clipboard on the right.
 * @param {string} titleText
 * @param {{ titleClass?: string, restoreLabel?: string, getSnippet: () => string }|null} [copy]
 */
function renderPanelHeader(titleText, copy = null) {
  const header = document.createElement('div');
  header.className = 'carousel-panel-header';

  const title = document.createElement('span');
  title.className = copy?.titleClass ?? 'carousel-color-panel-title';
  title.textContent = titleText;
  header.appendChild(title);

  if (copy) {
    const restoreLabel = copy.restoreLabel ?? 'Copy JSON';
    const copyBtn = createCopyButton({
      className: 'carousel-button carousel-button--compact carousel-section-copy',
      label: restoreLabel,
    });
    copyBtn.addEventListener('click', () => {
      copyText(copyBtn, copy.getSnippet(), restoreLabel);
    });
    header.appendChild(copyBtn);
  }

  return header;
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {CarouselDeck} deck
 * @param {PreviewSlot[]} previewSlots
 */
function renderMarginsPanel(theme, deck, previewSlots) {
  const panel = document.createElement('div');
  panel.className = 'carousel-debug-margins';

  panel.appendChild(renderPanelHeader('Margins (1080px)', {
    titleClass: 'carousel-debug-lineheights-title',
    restoreLabel: 'Copy margins JSON',
    getSnippet: () => marginsJsonSnippet(deck, theme),
  }));

  const fields = document.createElement('div');
  fields.className = 'carousel-lineheight-fields';
  colorControlRefs.marginInputs = [];

  /** @type {readonly { axis: 'horizontal' | 'vertical', label: string }[]} */
  const entries = [
    { axis: 'horizontal', label: 'Left / right' },
    { axis: 'vertical', label: 'Top / bottom' },
  ];

  for (const entry of entries) {
    const field = document.createElement('div');
    field.className = 'carousel-lineheight-field';

    const label = document.createElement('span');
    label.className = 'carousel-lineheight-label';
    label.textContent = entry.label;

    const inputWrap = document.createElement('div');
    inputWrap.className = 'carousel-lineheight-input-wrap';

    const input = document.createElement('input');
    input.type = 'number';
    input.className = 'carousel-lineheight-input';
    input.min = String(MARGIN_PERCENT_MIN);
    input.max = String(MARGIN_PERCENT_MAX);
    input.step = '1';
    input.dataset.marginAxis = entry.axis;
    input.value = String(marginPercentFromTheme(theme, entry.axis));
    input.setAttribute('aria-label', `${entry.label} margin percent`);

    const suffix = document.createElement('span');
    suffix.className = 'carousel-lineheight-suffix';
    suffix.textContent = '%';
    suffix.setAttribute('aria-hidden', 'true');

    input.addEventListener('input', () => {
      const parsed = Number(input.value);
      if (!Number.isFinite(parsed)) return;
      setDeckMarginPercent(deck, theme, entry.axis, parsed);
      refreshAllPreviews(previewSlots);
      studioStatePersistHandler?.();
    });

    input.addEventListener('change', () => {
      const parsed = Number(input.value);
      const next = Number.isFinite(parsed)
        ? parsed
        : marginPercentFromTheme(theme, entry.axis);
      setDeckMarginPercent(deck, theme, entry.axis, next);
      input.value = String(marginPercentFromTheme(theme, entry.axis));
      refreshAllPreviews(previewSlots);
      studioStatePersistHandler?.();
    });

    input._refreshMargin = () => {
      input.value = String(marginPercentFromTheme(theme, entry.axis));
    };

    inputWrap.append(input, suffix);
    field.append(label, inputWrap);
    fields.appendChild(field);
    colorControlRefs.marginInputs.push(input);
  }
  panel.appendChild(fields);

  const hint = document.createElement('p');
  hint.className = 'carousel-debug-lineheights-hint';
  hint.textContent = '0% = full bleed. Max 16%. Writes marginHorizontal / marginVertical in carousel.json.';
  panel.appendChild(hint);

  return panel;
}

/** @param {CarouselDeck} deck @param {import('./theme.js').DeckTheme} theme */
function marginsJsonSnippet(deck, theme) {
  const spec = deck.deck || {};
  /** @type {Record<string, string|number>} */
  const payload = {};
  const h = spec.marginHorizontal ?? spec.marginX;
  const v = spec.marginVertical ?? spec.marginY;
  if (h != null && h !== '') payload.marginHorizontal = h;
  if (v != null && v !== '') payload.marginVertical = v;
  if (Object.keys(payload).length === 0 && spec.margin != null && spec.margin !== '') {
    return JSON.stringify({ margin: spec.margin }, null, 2);
  }
  if (Object.keys(payload).length === 0) {
    payload.marginHorizontal = `${marginPercentFromTheme(theme, 'horizontal')}%`;
    payload.marginVertical = `${marginPercentFromTheme(theme, 'vertical')}%`;
  }
  return JSON.stringify(payload, null, 2);
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {CarouselDeck} deck
 * @param {PreviewSlot[]} previewSlots
 */
function renderLineHeightsPanel(theme, deck, previewSlots) {
  const panel = document.createElement('div');
  panel.className = 'carousel-debug-lineheights';

  panel.appendChild(renderPanelHeader('Line heights', {
    titleClass: 'carousel-debug-lineheights-title',
    restoreLabel: 'Copy lineHeights JSON',
    getSnippet: () => lineHeightsJsonSnippet(theme),
  }));

  const fields = document.createElement('div');
  fields.className = 'carousel-lineheight-fields';
  colorControlRefs.lineHeightInputs = [];

  for (const entry of LINE_HEIGHT_CONTROL_KEYS) {
    const field = document.createElement('div');
    field.className = 'carousel-lineheight-field';

    const label = document.createElement('span');
    label.className = 'carousel-lineheight-label';
    label.textContent = entry.label;

    const input = document.createElement('input');
    input.type = 'number';
    input.className = 'carousel-lineheight-input';
    input.min = '0.25';
    input.max = '4';
    input.step = '0.05';
    input.dataset.lineHeightKey = entry.key;
    input.value = String(lineHeightValueFor(theme, entry.key));
    input.setAttribute('aria-label', `${entry.label} line height`);

    input.addEventListener('input', () => {
      const parsed = Number(input.value);
      if (!Number.isFinite(parsed)) return;
      setDeckLineHeight(theme, deck, entry.key, parsed);
      refreshAllPreviews(previewSlots);
      studioStatePersistHandler?.();
    });

    input.addEventListener('change', () => {
      const parsed = Number(input.value);
      const next = Number.isFinite(parsed) ? parsed : 1;
      setDeckLineHeight(theme, deck, entry.key, next);
      input.value = String(lineHeightValueFor(theme, entry.key));
      refreshAllPreviews(previewSlots);
      studioStatePersistHandler?.();
    });

    input._refreshLineHeight = () => {
      input.value = String(lineHeightValueFor(theme, entry.key));
    };

    field.append(label, input);
    fields.appendChild(field);
    colorControlRefs.lineHeightInputs.push(input);
  }
  panel.appendChild(fields);

  const hint = document.createElement('p');
  hint.className = 'carousel-debug-lineheights-hint';
  hint.textContent = 'Blue slot = orange band × value. 1 = flush; 2 = double.';
  panel.appendChild(hint);

  return panel;
}

/** @param {import('./theme.js').DeckTheme} theme */
function lineHeightsJsonSnippet(theme) {
  /** @type {Record<string, number>} */
  const lineHeights = {};
  for (const { key } of LINE_HEIGHT_CONTROL_KEYS) {
    lineHeights[key] = lineHeightValueFor(theme, key);
  }
  return JSON.stringify({ lineHeights }, null, 2);
}

/** @param {CarouselDeck} deck */
function deckHasMotifStrip(deck) {
  const spec = deck.deck?.motifStrip;
  if (!spec || typeof spec !== 'object' || Array.isArray(spec)) return false;
  return typeof spec.src === 'string' && spec.src.trim().length > 0;
}

/** @param {CarouselDeck} deck */
function motifStripEnabledInDeck(deck) {
  if (!deckHasMotifStrip(deck)) return false;
  return /** @type {Record<string, unknown>} */ (deck.deck.motifStrip).enabled !== false;
}

/** @param {CarouselDeck} deck @param {boolean} enabled */
function setMotifStripEnabled(deck, enabled) {
  if (!deckHasMotifStrip(deck)) return;
  if (!deck.deck) deck.deck = {};
  const spec = deck.deck.motifStrip;
  if (!spec || typeof spec !== 'object' || Array.isArray(spec)) return;
  /** @type {Record<string, unknown>} */ (spec).enabled = enabled;
}

/** Horizontal pan: ±12 slide widths so a panorama can be shifted by whole slides. */
const MOTIF_OFFSET_X_MAX = CAROUSEL_SLIDE_WIDTH_PX * 12;
/** Vertical nudge: ± one slide (the strip is a short band). */
const MOTIF_OFFSET_Y_MAX = CAROUSEL_SLIDE_WIDTH_PX;
const MOTIF_SCALE_PERCENT_MIN = 20;
const MOTIF_SCALE_PERCENT_MAX = 200;
const MOTIF_SCALE_DELTA_MIN = MOTIF_SCALE_PERCENT_MIN - 100;
const MOTIF_SCALE_DELTA_MAX = MOTIF_SCALE_PERCENT_MAX - 100;

/** @param {CarouselDeck} deck */
function motifStripSpec(deck) {
  if (!deckHasMotifStrip(deck)) return null;
  return /** @type {Record<string, unknown>} */ (deck.deck.motifStrip);
}

/** @param {CarouselDeck} deck */
function motifScalePercent(deck) {
  const spec = motifStripSpec(deck);
  const raw = spec?.bandWidth;
  if (raw == null || raw === '') return 100;
  const str = String(raw).trim();
  const pct = /^(\d+(?:\.\d+)?)\s*%$/.exec(str);
  if (pct) {
    return Math.max(
      MOTIF_SCALE_PERCENT_MIN,
      Math.min(MOTIF_SCALE_PERCENT_MAX, Math.round(Number(pct[1]))),
    );
  }
  const num = Number(str);
  if (Number.isFinite(num) && num > 0 && num <= MOTIF_SCALE_PERCENT_MAX) {
    return Math.max(MOTIF_SCALE_PERCENT_MIN, Math.round(num));
  }
  return 100;
}

/** @param {CarouselDeck} deck */
function motifScaleDelta(deck) {
  return motifScalePercent(deck) - 100;
}

/** @param {number} delta */
function formatMotifScaleDelta(delta) {
  const n = Math.round(delta);
  return n > 0 ? `+${n}` : String(n);
}

/** @param {unknown} raw */
function parseMotifScaleDelta(raw) {
  const str = String(raw ?? '').trim().replace(/%/g, '');
  if (str === '' || str === '+' || str === '-' || str === '−') return null;
  const parsed = Number(str.replace('−', '-'));
  if (!Number.isFinite(parsed)) return null;
  return parsed;
}

/** @param {CarouselDeck} deck @param {number} percent */
function setMotifStripScale(deck, percent) {
  const spec = motifStripSpec(deck);
  if (!spec) return;
  const clamped = Math.round(Math.max(
    MOTIF_SCALE_PERCENT_MIN,
    Math.min(MOTIF_SCALE_PERCENT_MAX, percent),
  ));
  spec.bandWidth = `${clamped}%`;
}

/** @param {CarouselDeck} deck @param {number} delta */
function setMotifStripScaleDelta(deck, delta) {
  const clamped = Math.round(Math.max(
    MOTIF_SCALE_DELTA_MIN,
    Math.min(MOTIF_SCALE_DELTA_MAX, delta),
  ));
  setMotifStripScale(deck, 100 + clamped);
}

/** @param {CarouselDeck} deck @param {'offsetX' | 'offsetY'} key */
function motifOffsetValue(deck, key) {
  const spec = motifStripSpec(deck);
  if (!spec) return 0;
  const raw = spec[key];
  return Number.isFinite(Number(raw)) ? Number(raw) : 0;
}

/** @param {'offsetX' | 'offsetY'} key */
function motifOffsetMax(key) {
  return key === 'offsetX' ? MOTIF_OFFSET_X_MAX : MOTIF_OFFSET_Y_MAX;
}

/** @param {CarouselDeck} deck @param {'offsetX' | 'offsetY'} key @param {number} value */
function setMotifStripOffset(deck, key, value) {
  const spec = motifStripSpec(deck);
  if (!spec) return;
  const max = motifOffsetMax(key);
  const clamped = Math.round(Math.max(-max, Math.min(max, value)));
  if (clamped === 0) {
    delete spec[key];
  } else {
    spec[key] = clamped;
  }
}

/** @param {CarouselDeck} deck */
function deckHasPostCta(deck) {
  if (deck.deck?.cta && typeof deck.deck.cta === 'object') return true;
  for (const slide of deck.slides ?? []) {
    for (const variant of slide.variants ?? []) {
      if (variant.archetype === 'post_cta') return true;
    }
  }
  return false;
}

/** @param {CarouselDeck} deck */
function ensureDeckCta(deck) {
  if (!deck.deck) deck.deck = {};
  if (!deck.deck.cta || typeof deck.deck.cta !== 'object') {
    deck.deck.cta = {};
  }
  return deck.deck.cta;
}

function ctaQrSizeRaw(cta) {
  const nested = cta.qr && typeof cta.qr === 'object' && !Array.isArray(cta.qr)
    ? /** @type {Record<string, unknown>} */ (cta.qr).size
    : undefined;
  return nested;
}

/** @param {Record<string, unknown>} cta */
function ctaQrObject(cta) {
  /** @type {Record<string, unknown>} */
  const qr = (cta.qr && typeof cta.qr === 'object' && !Array.isArray(cta.qr))
    ? { .../** @type {Record<string, unknown>} */ (cta.qr) }
    : {};
  return qr;
}

/**
 * @param {Record<string, unknown>} cta
 * @param {keyof typeof CTA_LAYOUT_DEFAULTS} key
 */
function ctaLayoutValue(cta, key) {
  const fallback = CTA_LAYOUT_DEFAULTS[key];
  if (key === 'qrSize') return parseQrSizePercent(ctaQrSizeRaw(cta));
  const raw = cta[key];
  if (!Number.isFinite(raw)) return fallback;
  return /** @type {number} */ (raw);
}

/** @param {number} percent */
function formatQrSizePercent(percent) {
  return `${Math.round(percent)}%`;
}

/**
 * @param {CarouselDeck} deck
 * @param {keyof typeof CTA_LAYOUT_DEFAULTS} key
 * @param {number} value
 */
function setCtaLayoutValue(deck, key, value) {
  const cta = ensureDeckCta(deck);
  const { min, max } = CTA_LAYOUT_LIMITS[key];
  const clamped = Math.round(Math.max(min, Math.min(max, value)));
  if (key === 'qrSize') {
    const qr = ctaQrObject(cta);
    qr.size = formatQrSizePercent(clamped);
    cta.qr = qr;
    delete cta.qrSize;
    return;
  }
  cta[key] = clamped;
}

/**
 * @param {CarouselDeck} deck
 * @param {PreviewSlot[]} previewSlots
 */
function renderCtaSettingsPanel(deck, previewSlots) {
  const cta = ensureDeckCta(deck);
  const panel = document.createElement('div');
  panel.className = 'carousel-debug-cta';

  panel.appendChild(renderPanelHeader('CTA slide (1080px)', {
    titleClass: 'carousel-debug-lineheights-title',
    restoreLabel: 'Copy CTA layout JSON',
    getSnippet: () => ctaLayoutJsonSnippet(deck),
  }));

  const fields = document.createElement('div');
  fields.className = 'carousel-lineheight-fields';
  colorControlRefs.ctaInputs = [];

  /** @type {(keyof typeof CTA_LAYOUT_DEFAULTS)[]} */
  const entryKeys = ['featuredMaxHeight', 'qrSize', 'brandMaxHeight'];

  for (const key of entryKeys) {
    const entry = { key, ...CTA_LAYOUT_LIMITS[key] };
    const field = document.createElement('div');
    field.className = 'carousel-lineheight-field';

    const label = document.createElement('span');
    label.className = 'carousel-lineheight-label';
    label.textContent = entry.label;

    const input = document.createElement('input');
    input.type = 'number';
    input.className = 'carousel-lineheight-input';
    input.min = String(entry.min);
    input.max = String(entry.max);
    input.step = String(entry.step);
    input.dataset.ctaKey = entry.key;
    input.value = String(ctaLayoutValue(cta, entry.key));
    input.setAttribute(
      'aria-label',
      entry.key === 'qrSize' ? `${entry.label} percent of URL–footer slot` : `${entry.label} at 1080 canvas`,
    );

    let inputWrap = null;
    if (entry.key === 'qrSize') {
      inputWrap = document.createElement('div');
      inputWrap.className = 'carousel-lineheight-input-wrap';
      const suffix = document.createElement('span');
      suffix.className = 'carousel-lineheight-suffix';
      suffix.textContent = '%';
      suffix.setAttribute('aria-hidden', 'true');
      inputWrap.append(input, suffix);
    }

    input.addEventListener('input', () => {
      const parsed = Number(input.value);
      if (!Number.isFinite(parsed)) return;
      setCtaLayoutValue(deck, entry.key, parsed);
      refreshAllPreviews(previewSlots);
      studioStatePersistHandler?.();
    });

    input.addEventListener('change', () => {
      const parsed = Number(input.value);
      const next = Number.isFinite(parsed) ? parsed : CTA_LAYOUT_DEFAULTS[entry.key];
      setCtaLayoutValue(deck, entry.key, next);
      input.value = String(ctaLayoutValue(ensureDeckCta(deck), entry.key));
      refreshAllPreviews(previewSlots);
      studioStatePersistHandler?.();
    });

    field.append(label, inputWrap ?? input);
    fields.appendChild(field);
    colorControlRefs.ctaInputs.push(input);
  }
  panel.appendChild(fields);

  const hint = document.createElement('p');
  hint.className = 'carousel-debug-lineheights-hint';
  hint.textContent = 'Writes to deck.cta in carousel.json. Refresh preview after editing JSON by hand.';
  panel.appendChild(hint);

  return panel;
}

/** @param {CarouselDeck} deck */
function ctaLayoutJsonSnippet(deck) {
  const cta = deck.deck?.cta;
  if (!cta) {
    return JSON.stringify({
      cta: {
        featuredMaxHeight: CTA_LAYOUT_DEFAULTS.featuredMaxHeight,
        qr: { size: formatQrSizePercent(CTA_LAYOUT_DEFAULTS.qrSize) },
        brandMaxHeight: CTA_LAYOUT_DEFAULTS.brandMaxHeight,
      },
    }, null, 2);
  }
  const qr = ctaQrObject(cta);
  qr.size = formatQrSizePercent(ctaLayoutValue(cta, 'qrSize'));
  return JSON.stringify({
    cta: {
      featuredMaxHeight: ctaLayoutValue(cta, 'featuredMaxHeight'),
      qr,
      brandMaxHeight: ctaLayoutValue(cta, 'brandMaxHeight'),
    },
  }, null, 2);
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreviewSlot[]} previewSlots
 */
function renderPalettePanel(theme, previewSlots) {
  const panel = document.createElement('div');
  panel.className = 'carousel-color-panel carousel-color-panel--palette';

  panel.appendChild(renderPanelHeader('Text palette', {
    titleClass: 'carousel-float-section-title',
    restoreLabel: 'Copy palette JSON',
    getSnippet: () => paletteJsonSnippet(theme),
  }));

  const grid = document.createElement('div');
  grid.className = 'carousel-palette-grid';
  grid.setAttribute('role', 'list');
  colorControlRefs.paletteCards = [];

  for (const palette of CAROUSEL_PALETTES) {
    grid.appendChild(renderPaletteCard(palette, theme, previewSlots));
  }
  panel.appendChild(grid);

  const fields = document.createElement('div');
  fields.className = 'carousel-palette-chips';
  colorControlRefs.colorPickers = [];
  for (const entry of paletteBaseColorEntries(deckPaletteFromTheme(theme))) {
    const field = createPaletteColorField(theme, entry, previewSlots);
    fields.appendChild(field);
    const input = field.querySelector('input[type="color"]');
    if (input instanceof HTMLInputElement) {
      colorControlRefs.colorPickers.push(input);
    }
  }
  panel.appendChild(fields);

  syncThemeUi(theme, previewSlots);
  return panel;
}

/** @param {import('./theme.js').DeckTheme} theme */
function backgroundHueValueFor(theme) {
  ensureBackgroundWaveTheme(theme);
  const raw = theme.backgroundWave?.hueShift;
  const limits = BACKGROUND_WAVE_LIMITS.hueShift;
  if (!Number.isFinite(raw)) return limits.default ?? 0;
  return clampWaveValue(raw, limits.min, limits.max);
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {number} value
 */
function setBackgroundHueShift(theme, value) {
  ensureBackgroundWaveTheme(theme);
  const limits = BACKGROUND_WAVE_LIMITS.hueShift;
  const next = clampWaveValue(value, limits.min, limits.max);
  theme.backgroundWave = {
    ...theme.backgroundWave,
    hueShift: next,
  };
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreviewSlot[]} previewSlots
 * @returns {HTMLDivElement}
 */
function createHueShiftField(theme, previewSlots) {
  colorControlRefs.hueInputs = [];

  const field = document.createElement('div');
  field.className = 'carousel-wave-modifier-field carousel-wave-modifier-field--hue';

  const label = document.createElement('span');
  label.className = 'carousel-wave-modifier-label';
  label.textContent = 'Hue';

  const control = document.createElement('div');
  control.className = 'carousel-wave-modifier-control';

  const limits = BACKGROUND_WAVE_LIMITS.hueShift;
  const range = document.createElement('input');
  range.type = 'range';
  range.className = 'carousel-wave-modifier-range';
  range.min = String(limits.min);
  range.max = String(limits.max);
  range.step = String(limits.step);
  range.dataset.hueControl = 'range';
  range.value = String(backgroundHueValueFor(theme));
  range.setAttribute('aria-label', 'Background hue shift');

  const number = document.createElement('input');
  number.type = 'number';
  number.className = 'carousel-wave-modifier-number carousel-background-hue-number';
  number.min = String(limits.min);
  number.max = String(limits.max);
  number.step = String(limits.step);
  number.dataset.hueControl = 'number';
  number.value = String(backgroundHueValueFor(theme));
  number.setAttribute('aria-label', 'Background hue shift degrees');

  const applyHue = (value) => {
    if (!Number.isFinite(value)) return;
    setBackgroundHueShift(theme, value);
    refreshBackgroundUi(theme, previewSlots);
  };

  range.addEventListener('input', () => {
    const parsed = Number(range.value);
    if (!Number.isFinite(parsed)) return;
    number.value = String(parsed);
    applyHue(parsed);
  });

  number.addEventListener('input', () => {
    const parsed = Number(number.value);
    if (!Number.isFinite(parsed)) return;
    range.value = String(clampWaveValue(parsed, limits.min, limits.max));
    applyHue(parsed);
  });

  number.addEventListener('change', () => {
    const parsed = Number(number.value);
    const next = Number.isFinite(parsed) ? parsed : backgroundHueValueFor(theme);
    applyHue(next);
  });

  const refreshPair = () => {
    const value = backgroundHueValueFor(theme);
    range.value = String(value);
    number.value = String(value);
  };
  range._refreshHueInput = refreshPair;
  number._refreshHueInput = refreshPair;

  control.append(range, number);
  field.append(label, control);
  colorControlRefs.hueInputs.push(range, number);
  return field;
}

/**
 * @param {number} value
 * @param {number} min
 * @param {number} max
 */
function clampWaveValue(value, min, max) {
  return Math.max(min, Math.min(max, value));
}

/** @param {import('./theme.js').DeckTheme} theme @param {'intensity' | 'color' | 'variety' | 'blur' | 'radius' | 'phase'} key */
function waveScalarValueFor(theme, key) {
  ensureBackgroundWaveTheme(theme);
  const wave = theme.backgroundWave;
  if (!wave) return BACKGROUND_WAVE_LIMITS[key].default ?? 0;
  const raw = wave[key];
  const limits = BACKGROUND_WAVE_LIMITS[key];
  if (!Number.isFinite(raw)) return limits.default ?? 0;
  return clampWaveValue(raw, limits.min, limits.max);
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {'intensity' | 'color' | 'variety' | 'blur' | 'radius' | 'phase'} key
 * @param {number} value
 */
function setWaveScalar(theme, key, value) {
  ensureBackgroundWaveTheme(theme);
  const limits = BACKGROUND_WAVE_LIMITS[key];
  const next = clampWaveValue(value, limits.min, limits.max);
  theme.backgroundWave = {
    ...theme.backgroundWave,
    [key]: next,
  };
}

/** @param {import('./theme.js').DeckTheme} theme */
function waveLobesValueFor(theme) {
  ensureBackgroundWaveTheme(theme);
  const lobes = theme.backgroundWave?.lobes;
  return Number.isFinite(lobes) ? lobes : null;
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {number|null} value
 */
function setWaveLobes(theme, value) {
  ensureBackgroundWaveTheme(theme);
  let lobes = null;
  if (value != null && Number.isFinite(value)) {
    lobes = clampWaveValue(
      Math.round(value),
      BACKGROUND_WAVE_LIMITS.lobes.min,
      BACKGROUND_WAVE_LIMITS.lobes.max,
    );
  }
  theme.backgroundWave = {
    ...theme.backgroundWave,
    lobes,
  };
}

/** @param {string} styleId */
function waveLobesHintText(styleId) {
  if (styleId === 'mesh-corners') {
    return 'Bottom travelers. Leave empty for auto (2–4).';
  }
  return 'Mid-strip pools. Leave empty for auto (slide count + 2).';
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreviewSlot[]} previewSlots
 */
function renderWaveModifierPanel(theme, previewSlots) {
  const panel = document.createElement('div');
  panel.className = 'carousel-wave-controls';
  colorControlRefs.waveModifierPanel = panel;
  colorControlRefs.waveInputs = [];

  const fields = document.createElement('div');
  fields.className = 'carousel-wave-modifier-fields';

  fields.appendChild(createHueShiftField(theme, previewSlots));

  const geometryFields = document.createElement('div');
  geometryFields.className = 'carousel-wave-modifier-fields carousel-wave-modifier-fields--geometry';
  colorControlRefs.waveGeometryFields = geometryFields;

  for (const entry of WAVE_SCALAR_CONTROL_KEYS) {
    const limits = BACKGROUND_WAVE_LIMITS[entry.key];
    const field = document.createElement('div');
    field.className = 'carousel-wave-modifier-field';

    const label = document.createElement('span');
    label.className = 'carousel-wave-modifier-label';
    label.textContent = entry.label;

    const control = document.createElement('div');
    control.className = 'carousel-wave-modifier-control';

    const range = document.createElement('input');
    range.type = 'range';
    range.className = 'carousel-wave-modifier-range';
    range.min = String(limits.min);
    range.max = String(limits.max);
    range.step = String(limits.step);
    range.dataset.waveKey = entry.key;
    range.dataset.waveControl = 'range';
    range.value = String(waveScalarValueFor(theme, entry.key));
    range.setAttribute('aria-label', `${entry.label} wave modifier`);

    const number = document.createElement('input');
    number.type = 'number';
    number.className = 'carousel-wave-modifier-number';
    number.min = String(limits.min);
    number.max = String(limits.max);
    number.step = String(limits.step);
    number.dataset.waveKey = entry.key;
    number.dataset.waveControl = 'number';
    number.value = String(waveScalarValueFor(theme, entry.key));
    number.setAttribute('aria-label', `${entry.label} wave modifier value`);

    const applyScalar = (value) => {
      if (!Number.isFinite(value)) return;
      setWaveScalar(theme, entry.key, value);
      refreshBackgroundUi(theme, previewSlots);
    };

    range.addEventListener('input', () => {
      const parsed = Number(range.value);
      if (!Number.isFinite(parsed)) return;
      number.value = String(parsed);
      applyScalar(parsed);
    });

    number.addEventListener('input', () => {
      const parsed = Number(number.value);
      if (!Number.isFinite(parsed)) return;
      range.value = String(clampWaveValue(parsed, limits.min, limits.max));
      applyScalar(parsed);
    });

    number.addEventListener('change', () => {
      const parsed = Number(number.value);
      const next = Number.isFinite(parsed) ? parsed : waveScalarValueFor(theme, entry.key);
      applyScalar(next);
    });

    const refreshPair = () => {
      const value = waveScalarValueFor(theme, entry.key);
      range.value = String(value);
      number.value = String(value);
    };
    range._refreshWaveInput = refreshPair;
    number._refreshWaveInput = refreshPair;

    control.append(range, number);
    field.append(label, control);
    geometryFields.appendChild(field);
    colorControlRefs.waveInputs.push(range, number);
  }

  const lobesField = document.createElement('div');
  lobesField.className = 'carousel-wave-modifier-field carousel-wave-modifier-field--lobes';

  const lobesLabel = document.createElement('span');
  lobesLabel.className = 'carousel-wave-modifier-label';
  lobesLabel.textContent = 'Lobes';

  const lobesWrap = document.createElement('div');
  lobesWrap.className = 'carousel-wave-modifier-control';

  const lobesInput = document.createElement('input');
  lobesInput.type = 'number';
  lobesInput.className = 'carousel-wave-modifier-number carousel-wave-modifier-number--lobes';
  lobesInput.min = String(BACKGROUND_WAVE_LIMITS.lobes.min);
  lobesInput.max = String(BACKGROUND_WAVE_LIMITS.lobes.max);
  lobesInput.step = String(BACKGROUND_WAVE_LIMITS.lobes.step);
  lobesInput.placeholder = 'Auto';
  lobesInput.dataset.waveKey = 'lobes';
  lobesInput.dataset.waveControl = 'lobes';
  const lobesValue = waveLobesValueFor(theme);
  lobesInput.value = lobesValue == null ? '' : String(lobesValue);
  lobesInput.setAttribute('aria-label', 'Lobes wave modifier');

  lobesInput.addEventListener('input', () => {
    const raw = lobesInput.value.trim();
    if (raw === '') {
      setWaveLobes(theme, null);
      refreshBackgroundUi(theme, previewSlots);
      return;
    }
    const parsed = Number(raw);
    if (!Number.isFinite(parsed)) return;
    setWaveLobes(theme, parsed);
    refreshBackgroundUi(theme, previewSlots);
  });

  lobesInput.addEventListener('change', () => {
    const raw = lobesInput.value.trim();
    if (raw === '') {
      setWaveLobes(theme, null);
    } else {
      const parsed = Number(raw);
      setWaveLobes(theme, Number.isFinite(parsed) ? parsed : null);
    }
    const next = waveLobesValueFor(theme);
    lobesInput.value = next == null ? '' : String(next);
    refreshBackgroundUi(theme, previewSlots);
  });

  lobesInput._refreshWaveInput = () => {
    const next = waveLobesValueFor(theme);
    lobesInput.value = next == null ? '' : String(next);
  };

  const lobesHint = document.createElement('span');
  lobesHint.className = 'carousel-wave-lobes-hint';
  lobesHint.textContent = waveLobesHintText(theme.backgroundWave?.style ?? 'drift');

  lobesWrap.append(lobesInput, lobesHint);
  lobesField.append(lobesLabel, lobesWrap);
  geometryFields.appendChild(lobesField);
  colorControlRefs.waveInputs.push(lobesInput);

  fields.appendChild(geometryFields);
  panel.appendChild(fields);
  return panel;
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreviewSlot[]} previewSlots
 */
function renderWavePaletteSection(theme, previewSlots) {
  const section = document.createElement('div');
  section.className = 'carousel-wave-palette-section';
  colorControlRefs.wavePaletteSection = section;

  const linkRow = document.createElement('label');
  linkRow.className = 'carousel-wave-palette-link';

  const linkToggle = document.createElement('input');
  linkToggle.type = 'checkbox';
  linkToggle.className = 'carousel-wave-palette-link-toggle';
  linkToggle.checked = !isWavePaletteLinked(theme);
  linkToggle.setAttribute('aria-label', 'Separate wave palette from text palette');
  colorControlRefs.wavePaletteLinkToggle = linkToggle;

  const linkLabel = document.createElement('span');
  linkLabel.textContent = 'Separate wave palette';

  linkRow.append(linkToggle, linkLabel);
  linkToggle.addEventListener('change', () => {
    setWavePaletteLinked(theme, !linkToggle.checked, previewSlots);
  });

  const linkedHint = document.createElement('p');
  linkedHint.className = 'carousel-wave-palette-linked-hint';
  linkedHint.textContent = 'Wash uses text palette Base, Accents, and Muted.';

  const fields = document.createElement('div');
  fields.className = 'carousel-palette-chips carousel-wave-palette-chips';
  colorControlRefs.waveColorPickers = [];

  const waveColors = deckWavePaletteFromTheme(theme)
    ?? normalizeWavePalette({}, panoramaPaletteFromTextPalette(deckPaletteFromTheme(theme)));
  for (const entry of wavePaletteColorEntries(waveColors)) {
    const field = createWavePaletteColorField(theme, entry, previewSlots);
    fields.appendChild(field);
    const input = field.querySelector('input[type="color"]');
    if (input instanceof HTMLInputElement) {
      colorControlRefs.waveColorPickers.push(input);
    }
  }

  section.append(linkRow, linkedHint, fields);
  syncWavePaletteSectionUi(theme);
  return section;
}

/** @param {import('./theme.js').DeckTheme} theme */
function syncWavePaletteSectionUi(theme) {
  const linked = isWavePaletteLinked(theme);
  const section = colorControlRefs.wavePaletteSection;
  if (section) {
    section.classList.toggle('carousel-wave-palette-section--linked', linked);
  }
  const toggle = colorControlRefs.wavePaletteLinkToggle;
  if (toggle) {
    toggle.checked = !linked;
  }
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {boolean} linked
 * @param {PreviewSlot[]} previewSlots
 */
function setWavePaletteLinked(theme, linked, previewSlots) {
  theme.wavePaletteLinked = linked;
  if (linked) {
    theme.wavePalette = null;
  } else {
    theme.wavePalette = normalizeWavePalette(
      theme.wavePalette && typeof theme.wavePalette === 'object' ? theme.wavePalette : {},
      panoramaPaletteFromTextPalette(deckPaletteFromTheme(theme)),
    );
  }
  syncWavePaletteDeckSpec(theme);
  refreshBackgroundUi(theme, previewSlots);
}

/** @param {import('./theme.js').DeckTheme} theme */
function syncWavePaletteDeckSpec(theme) {
  if (!studioDeckRef?.deck) return;
  const wavePalette = deckWavePaletteFromTheme(theme);
  if (wavePalette) {
    studioDeckRef.deck.wavePalette = { ...wavePalette };
    studioDeckRef.deck.wavePaletteLinked = false;
  } else {
    delete studioDeckRef.deck.wavePalette;
    studioDeckRef.deck.wavePaletteLinked = true;
  }
}

/** @param {import('./theme.js').DeckTheme} theme */
function ensureWavePaletteObject(theme) {
  if (!theme.wavePalette || typeof theme.wavePalette !== 'object') {
    theme.wavePalette = normalizeWavePalette(
      {},
      panoramaPaletteFromTextPalette(deckPaletteFromTheme(theme)),
    );
  }
  theme.wavePaletteLinked = false;
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {import('./palettes.js').WavePaletteColorEntry} entry
 * @param {PreviewSlot[]} previewSlots
 */
function createWavePaletteColorField(theme, entry, previewSlots) {
  const row = document.createElement('label');
  row.className = 'carousel-palette-color-chip carousel-wave-palette-color-chip';

  const label = document.createElement('span');
  label.className = 'carousel-color-label';
  label.textContent = entry.label;

  const input = document.createElement('input');
  input.type = 'color';
  input.className = 'carousel-color-input';
  input.dataset.waveColorKey = entry.key;
  input.value = waveThemeColorHex(theme, entry.key);
  input.setAttribute('aria-label', `Wave ${entry.label} color`);

  input.addEventListener('input', () => {
    setWaveThemeColor(theme, entry.key, input.value);
    syncWavePaletteDeckSpec(theme);
    refreshBackgroundUi(theme, previewSlots);
  });

  input._refreshWaveColorPicker = () => {
    input.value = waveThemeColorHex(theme, entry.key);
  };

  row.append(input, label);
  return row;
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {import('./palettes.js').WavePaletteColorEntry['key']} key
 */
function waveThemeColorHex(theme, key) {
  const wave = wavePaletteFromTheme(theme);
  if (key === 'accent2') {
    return toColorInputValue(wave.accent2 || wave.accent1);
  }
  return toColorInputValue(wave[key]);
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {import('./palettes.js').WavePaletteColorEntry['key']} key
 * @param {string} hex
 */
function setWaveThemeColor(theme, key, hex) {
  ensureWavePaletteObject(theme);
  const value = toColorInputValue(hex);
  if (!theme.wavePalette || typeof theme.wavePalette !== 'object') return;
  if (key === 'accent2') {
    theme.wavePalette.accent2 = value;
    return;
  }
  theme.wavePalette[key] = value;
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreviewSlot[]} previewSlots
 */
function renderGradientPanel(theme, previewSlots) {
  const panel = document.createElement('div');
  panel.className = 'carousel-color-panel carousel-color-panel--wave';

  panel.appendChild(renderPanelHeader('Background wave', {
    titleClass: 'carousel-float-section-title',
    restoreLabel: 'Copy wave JSON',
    getSnippet: () => waveJsonSnippet(theme),
  }));

  panel.appendChild(renderWavePaletteSection(theme, previewSlots));

  const grid = document.createElement('div');
  grid.className = 'carousel-gradient-grid';
  grid.setAttribute('role', 'list');
  colorControlRefs.gradientCards = [];

  ensureBackgroundWaveTheme(theme);

  for (const style of BACKGROUND_WAVE_STYLES) {
    grid.appendChild(renderWaveStyleCard(style, theme, previewSlots));
  }
  panel.appendChild(grid);
  panel.appendChild(renderWaveModifierPanel(theme, previewSlots));

  return panel;
}

/** @param {import('./theme.js').DeckTheme} theme */
function ensureBackgroundWaveTheme(theme) {
  theme.backgroundGradientMode = 'panoramic-wave';
  theme.backgroundGradientPreset = null;
  theme.backgroundGradient = 'solid';
  if (!theme.backgroundWave) {
    theme.backgroundWave = parseBackgroundWaveConfig({});
  }
}

/**
 * @param {{ id: string, label: string }} style
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreviewSlot[]} previewSlots
 */
function renderWaveStyleCard(style, theme, previewSlots) {
  const card = document.createElement('button');
  card.type = 'button';
  card.className = 'carousel-gradient-card';
  card.setAttribute('role', 'listitem');
  card.dataset.gradientId = style.id;
  card.title = style.label;

  const activeStyle = theme.backgroundWave?.style ?? 'drift';
  if (activeStyle === style.id) {
    card.classList.add('carousel-gradient-card--active');
  }

  const preview = document.createElement('span');
  preview.className = 'carousel-gradient-card-preview';
  preview.style.background = waveStylePreviewCss(theme, style.id);

  const name = document.createElement('span');
  name.className = 'carousel-gradient-card-name';
  const shortNames = { none: 'None', drift: 'Drift', 'mesh-corners': 'Mesh corners' };
  name.textContent = shortNames[style.id] ?? style.label;

  card.append(preview, name);
  card.addEventListener('click', () => {
    ensureBackgroundWaveTheme(theme);
    theme.backgroundWave = {
      ...theme.backgroundWave,
      style: style.id,
    };
    refreshBackgroundUi(theme, previewSlots);
  });

  colorControlRefs.gradientCards?.push(card);
  return card;
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {string} styleId
 */
function waveStylePreviewCss(theme, styleId) {
  ensureBackgroundWaveTheme(theme);
  const palette = paletteWithBackgroundHueShift(
    wavePaletteFromTheme(theme),
    theme.backgroundWave,
  );
  const base = palette.background;
  if (styleId === 'none') {
    return base;
  }
  const wave = theme.backgroundWave;
  const presence = waveVisualStrength(wave?.intensity ?? BACKGROUND_WAVE_LIMITS.intensity.default);
  const chroma = waveColorStrength(wave?.color);
  const extended = extendedBackgroundPalette(
    {
      background: base,
      accent1: palette.accent1,
      accent2: palette.accent2 || palette.accent1,
      muted: palette.muted,
    },
    wave?.variety,
  );
  const wash = (hex) => mixHex(base, hex, 0.04 + chroma * 0.62 + presence * 0.16);
  const n = extended.length;
  const at = (fraction) => {
    const t = Math.max(0, Math.min(1, fraction));
    return wash(extended[Math.round(t * Math.max(0, n - 1))]);
  };
  const accent1 = at(0.22);
  const accent2 = at(0.48);
  const muted = at(0.72);
  const washSpread = 66;
  if (styleId === 'mesh-corners') {
    return [
      `radial-gradient(ellipse 70% 90% at 8% 92%, ${accent1} 0%, transparent ${washSpread + 6}%)`,
      `radial-gradient(ellipse 70% 90% at 92% 90%, ${accent2} 0%, transparent ${washSpread + 4}%)`,
      `radial-gradient(ellipse 55% 70% at 90% 12%, ${muted} 0%, transparent ${washSpread}%)`,
      `radial-gradient(ellipse 80% 85% at 50% 48%, ${base} 0%, transparent ${washSpread + 10}%)`,
      base,
    ].join(', ');
  }
  return [
    `radial-gradient(ellipse 55% 80% at 18% 62%, ${accent1} 0%, transparent ${washSpread + 8}%)`,
    `radial-gradient(ellipse 45% 70% at 47% 38%, ${muted} 0%, transparent ${washSpread + 6}%)`,
    `radial-gradient(ellipse 50% 75% at 73% 72%, ${accent2} 0%, transparent ${washSpread + 4}%)`,
    `radial-gradient(ellipse 40% 60% at 88% 28%, ${accent1} 0%, transparent ${washSpread}%)`,
    base,
  ].join(', ');
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {import('./palettes.js').CarouselPalette} palette
 */
function applyPaletteTheme(theme, palette) {
  const p = palette.palette;
  theme.background = p.background;
  theme.text = p.text;
  theme.muted = p.muted;
  theme.accent1 = p.accent1;
  theme.accent2 = p.accent2 ?? null;
  rebuildBackgroundGradient(theme);
}

/**
 * @param {import('./palettes.js').CarouselPalette} palette
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreviewSlot[]} previewSlots
 */
function renderPaletteCard(palette, theme, previewSlots) {
  const card = document.createElement('button');
  card.type = 'button';
  card.className = 'carousel-palette-card';
  card.setAttribute('role', 'listitem');
  card.dataset.paletteId = palette.id;
  card.setAttribute('aria-label', palette.label);
  card.title = palette.label;
  if (resolvePaletteId(theme, studioPaletteId) === palette.id) {
    card.classList.add('carousel-palette-card--active');
  }

  const name = document.createElement('span');
  name.className = 'carousel-palette-card-name';
  name.textContent = palette.label;

  const dots = document.createElement('span');
  dots.className = 'carousel-palette-card-dots';
  dots.setAttribute('aria-hidden', 'true');
  for (const entry of paletteBaseColorEntries(palette.palette)) {
    const dot = document.createElement('span');
    dot.className = 'carousel-palette-mini-swatch';
    dot.style.setProperty('--swatch-color', entry.hex);
    if (isLightHex(entry.hex)) {
      dot.classList.add('carousel-palette-mini-swatch--light');
    } else {
      dot.style.background = entry.hex;
    }
    dot.title = entry.label;
    dots.appendChild(dot);
  }

  card.append(name, dots);
  card.addEventListener('click', () => {
    studioPaletteId = palette.id;
    applyPaletteTheme(theme, palette);
    syncThemeUi(theme, previewSlots);
  });

  colorControlRefs.paletteCards?.push(card);
  return card;
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {import('./palettes.js').PaletteColorEntry} entry
 * @param {PreviewSlot[]} previewSlots
 */
function createPaletteColorField(theme, entry, previewSlots) {
  const row = document.createElement('label');
  row.className = 'carousel-palette-color-chip';

  const label = document.createElement('span');
  label.className = 'carousel-color-label';
  label.textContent = entry.label;

  const input = document.createElement('input');
  input.type = 'color';
  input.className = 'carousel-color-input';
  input.dataset.colorKey = entry.key;
  input.value = themeColorHex(theme, entry.key);
  input.setAttribute('aria-label', `${entry.label} color`);

  input.addEventListener('input', () => {
    setThemeColor(theme, entry.key, input.value);
    rebuildBackgroundGradient(theme);
    studioPaletteId = resolvePaletteId(theme, studioPaletteId);
    syncThemeUi(theme, previewSlots);
  });

  input._refreshColorPicker = () => {
    input.value = themeColorHex(theme, entry.key);
  };

  row.append(input, label);
  return row;
}

/** @param {import('./theme.js').DeckTheme} theme @param {import('./palettes.js').PaletteColorEntry['key']} key */
function themeColorHex(theme, key) {
  if (key === 'accent2') {
    return toColorInputValue(theme.accent2 || theme.accent1);
  }
  return toColorInputValue(theme[key]);
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {import('./palettes.js').PaletteColorEntry['key']} key
 * @param {string} hex
 */
function setThemeColor(theme, key, hex) {
  const value = toColorInputValue(hex);
  if (key === 'accent2') {
    theme.accent2 = value;
    return;
  }
  theme[key] = value;
}

/** @param {import('./theme.js').DeckTheme} theme @param {PreviewSlot[]} previewSlots */
function syncThemeUi(theme, previewSlots) {
  studioPaletteId = resolvePaletteId(theme, studioPaletteId);
  const activePaletteId = studioPaletteId;
  for (const card of colorControlRefs.paletteCards ?? []) {
    card.classList.toggle('carousel-palette-card--active', card.dataset.paletteId === activePaletteId);
  }

  ensureBackgroundWaveTheme(theme);
  const activeStyle = theme.backgroundWave?.style ?? 'drift';
  for (const card of colorControlRefs.gradientCards ?? []) {
    const styleId = card.dataset.gradientId || 'drift';
    card.classList.toggle('carousel-gradient-card--active', styleId === activeStyle);
    const preview = card.querySelector('.carousel-gradient-card-preview');
    if (preview instanceof HTMLElement) {
      preview.style.background = waveStylePreviewCss(theme, styleId);
    }
  }

  for (const input of colorControlRefs.colorPickers ?? []) {
    if (typeof input._refreshColorPicker === 'function') {
      input._refreshColorPicker();
    }
  }

  syncWavePaletteSectionUi(theme);
  for (const input of colorControlRefs.waveColorPickers ?? []) {
    if (typeof input._refreshWaveColorPicker === 'function') {
      input._refreshWaveColorPicker();
    }
  }

  for (const input of colorControlRefs.lineHeightInputs ?? []) {
    if (typeof input._refreshLineHeight === 'function') {
      input._refreshLineHeight();
    }
  }

  for (const input of colorControlRefs.marginInputs ?? []) {
    if (typeof input._refreshMargin === 'function') {
      input._refreshMargin();
    }
  }

  for (const input of colorControlRefs.motifOffsetInputs ?? []) {
    if (typeof input._refreshMotifOffset === 'function') {
      input._refreshMotifOffset();
    }
  }

  for (const input of colorControlRefs.hueInputs ?? []) {
    if (typeof input._refreshHueInput === 'function') {
      input._refreshHueInput();
    }
  }

  const waveModifiersActive = activeStyle !== 'none';
  colorControlRefs.waveGeometryFields?.classList.toggle(
    'carousel-wave-modifier-fields--geometry--hidden',
    !waveModifiersActive,
  );
  for (const input of colorControlRefs.waveInputs ?? []) {
    input.disabled = !waveModifiersActive;
    if (typeof input._refreshWaveInput === 'function') {
      input._refreshWaveInput();
    }
  }
  const lobesHint = colorControlRefs.waveModifierPanel?.querySelector('.carousel-wave-lobes-hint');
  if (lobesHint instanceof HTMLElement) {
    lobesHint.textContent = waveLobesHintText(activeStyle);
  }

  refreshAllPreviews(previewSlots);
  studioStatePersistHandler?.();
}

/** @param {import('./theme.js').DeckTheme} theme */
function paletteJsonSnippet(theme) {
  return JSON.stringify({ palette: deckPaletteFromTheme(theme) }, null, 2);
}

/**
 * Paste-ready `backgroundGradient` + `backgroundWave` for `deck` in carousel.json.
 * @param {import('./theme.js').DeckTheme} theme
 */
function waveJsonSnippet(theme) {
  ensureBackgroundWaveTheme(theme);
  /** @type {Record<string, unknown>} */
  const payload = {
    backgroundGradient: 'solid',
    backgroundWave: compactBackgroundWaveForExport(theme.backgroundWave),
  };
  const wavePalette = deckWavePaletteFromTheme(theme);
  if (wavePalette) {
    payload.wavePalette = wavePalette;
  }
  return JSON.stringify(payload, null, 2);
}

/**
 * @param {import('./background-panorama.js').BackgroundWaveConfig} wave
 * @returns {import('./background-panorama.js').BackgroundWaveConfig}
 */
function compactBackgroundWaveForExport(wave) {
  if (wave.style === 'none') {
    if (wave.hueShift) {
      return { style: 'none', hueShift: Math.round(wave.hueShift) };
    }
    return { style: 'none' };
  }
  /** @type {import('./background-panorama.js').BackgroundWaveConfig} */
  const out = {
    style: wave.style,
    lobes: wave.lobes,
    intensity: roundWaveExportNumber(wave.intensity, 2),
    color: roundWaveExportNumber(wave.color ?? BACKGROUND_WAVE_LIMITS.color.default, 2),
    variety: roundWaveExportNumber(wave.variety ?? BACKGROUND_WAVE_LIMITS.variety.default, 2),
    blur: roundWaveExportNumber(wave.blur, 2),
    radius: roundWaveExportNumber(wave.radius ?? BACKGROUND_WAVE_LIMITS.radius.default, 2),
    phase: wave.phase,
    hueShift: wave.hueShift,
  };
  if (out.lobes == null) delete out.lobes;
  if (Math.abs(out.color - BACKGROUND_WAVE_LIMITS.color.default) < 0.001) delete out.color;
  if (Math.abs(out.variety - BACKGROUND_WAVE_LIMITS.variety.default) < 0.001) delete out.variety;
  if (Math.abs(out.radius - BACKGROUND_WAVE_LIMITS.radius.default) < 0.001) delete out.radius;
  if (!out.phase) delete out.phase;
  if (!out.hueShift) delete out.hueShift;
  return out;
}

/** @param {number} value @param {number} digits */
function roundWaveExportNumber(value, digits) {
  if (!Number.isFinite(value)) return value;
  const factor = 10 ** digits;
  return Math.round(value * factor) / factor;
}

function refreshBackgroundUi(theme, previewSlots) {
  ensureBackgroundWaveTheme(theme);
  if (studioDeckRef?.deck) {
    if (theme.backgroundWave) {
      studioDeckRef.deck.backgroundWave = { ...theme.backgroundWave };
    }
    syncWavePaletteDeckSpec(theme);
  }
  rebuildBackgroundGradient(theme);
  syncThemeUi(theme, previewSlots);
}

/** @param {import('./theme.js').DeckTheme} theme @param {PreviewSlot[]} previewSlots */
function refreshThemeColors(theme, previewSlots) {
  refreshBackgroundUi(theme, previewSlots);
}

/** @param {PreviewSlot[]} previewSlots */
async function refreshAllPreviews(previewSlots) {
  await Promise.all(previewSlots.map((slot) => slot.refresh()));
  studioVisionStripReschedule?.();
}

/** @param {PreviewSlot[]} previewSlots */
function attachPreviewResizeObservers(previewSlots) {
  for (const slot of previewSlots) {
    let resizeTimer = 0;
    const observer = new ResizeObserver(() => {
      window.clearTimeout(resizeTimer);
      resizeTimer = window.setTimeout(() => {
        slot.refresh();
      }, 80);
    });
    observer.observe(slot.previewFrame);
  }
}

/** @param {string|null|undefined} value */
function toColorInputValue(value) {
  if (!value || typeof value !== 'string') return '#d69a80';
  const trimmed = value.trim();
  if (/^#[0-9a-fA-F]{6}$/.test(trimmed)) return trimmed;
  if (/^[0-9a-fA-F]{6}$/.test(trimmed)) return `#${trimmed}`;
  return '#d69a80';
}

/** @param {string} hex */
function isLightHex(hex) {
  const match = /^#?([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})$/.exec(hex.trim());
  if (!match) return false;
  const channels = match.slice(1).map((part) => parseInt(part, 16) / 255);
  const luminance = 0.2126 * channels[0] + 0.7152 * channels[1] + 0.0722 * channels[2];
  return luminance > 0.72;
}

/** @param {string} value */
function slugify(value) {
  return value.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '');
}

/** @param {string} value */
function escapeHtml(value) {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;');
}
