/**
 * @typedef {Object} ChannelState
 * @property {string} status
 */

/**
 * @typedef {Object} CalendarEntry
 * @property {string} bundle
 * @property {string} section
 * @property {string} title
 * @property {string} [type]
 * @property {string[]} [categories]
 * @property {string} planned
 * @property {boolean} draft
 * @property {Record<string, ChannelState>} channels
 */

/**
 * @typedef {Object} CalendarData
 * @property {string} generatedAt
 * @property {string[]} channels
 * @property {CalendarEntry[]} entries
 */

const WEEKDAYS = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
const MONTHS = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December',
];

const CHANNEL_LABELS = {
  'site-en': 'English',
  'site-es': 'Spanish',
  linkedin: 'LinkedIn',
  facebook: 'Facebook',
  'substack-en': 'Substack',
  'substack-es': 'Substack ES',
};

const CHANNEL_SHORT = {
  'site-en': 'EN',
  'site-es': 'ES',
  linkedin: 'LI',
  facebook: 'FB',
  'substack-en': 'SB',
};

/** @param {string} id */
function channelLegendLabel(id) {
  const name = CHANNEL_LABELS[id] || id;
  const short = CHANNEL_SHORT[id];
  return short ? `${name} (${short})` : name;
}

/** @type {readonly string[]} */
const CHANNEL_BOX_ORDER = ['site-en', 'site-es', 'linkedin', 'facebook', 'substack-en'];

const STORAGE_KEY = 'publish-calendar-view';

/**
 * @param {string} dayKey
 * @returns {{ year: number, month: number, selectedDay: string } | null}
 */
function parseDayKeyState(dayKey) {
  const match = /^(\d{4})-(\d{2})-(\d{2})$/.exec(dayKey);
  if (!match) return null;
  const year = Number(match[1]);
  const monthIndex = Number(match[2]) - 1;
  const day = Number(match[3]);
  if (monthIndex < 0 || monthIndex > 11 || day < 1 || day > 31) return null;
  const date = new Date(year, monthIndex, day);
  if (
    date.getFullYear() !== year
    || date.getMonth() !== monthIndex
    || date.getDate() !== day
  ) {
    return null;
  }
  return { year, month: monthIndex, selectedDay: dayKey };
}

/** @returns {{ year: number, month: number, selectedDay: string } | null} */
function parseViewQuery(params) {
  const y = params.get('y');
  const m = params.get('m');
  if (!y || !m) return null;
  const year = Number(y);
  const month = Number(m) - 1;
  if (Number.isNaN(year) || Number.isNaN(month) || month < 0 || month > 11) return null;
  const d = params.get('d');
  const dayState = d ? parseDayKeyState(d) : null;
  return {
    year,
    month,
    selectedDay: dayState?.selectedDay ?? formatDayKey(new Date(year, month, 1)),
  };
}

/** @returns {{ year: number, month: number, selectedDay: string } | null} */
function readStorageState() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw);
    if (typeof parsed.year !== 'number' || typeof parsed.month !== 'number') return null;
    if (parsed.month < 0 || parsed.month > 11) return null;
    const selectedDay = typeof parsed.selectedDay === 'string'
      ? (parseDayKeyState(parsed.selectedDay)?.selectedDay ?? formatDayKey(new Date()))
      : formatDayKey(new Date());
    return { year: parsed.year, month: parsed.month, selectedDay };
  } catch {
    return null;
  }
}

/** @returns {{ year: number, month: number, selectedDay: string }} */
function readPersistedCalendarState() {
  const now = new Date();
  const fallback = {
    year: now.getFullYear(),
    month: now.getMonth(),
    selectedDay: formatDayKey(now),
  };

  const fromQuery = parseViewQuery(new URLSearchParams(window.location.search));
  if (fromQuery) return fromQuery;

  const fromStorage = readStorageState();
  if (fromStorage) return fromStorage;

  const fromHash = parseDayKeyState(window.location.hash.replace(/^#/, ''));
  if (fromHash) return fromHash;

  return fallback;
}

/** @param {{ year: number, month: number, selectedDay: string }} state */
function persistCalendarState(state) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify({
      year: state.year,
      month: state.month,
      selectedDay: state.selectedDay,
    }));
  } catch {
    // ignore quota / private mode
  }

  const params = new URLSearchParams();
  params.set('y', String(state.year));
  params.set('m', String(state.month + 1));
  params.set('d', state.selectedDay);
  const nextUrl = `${window.location.pathname}?${params.toString()}`;
  const currentUrl = `${window.location.pathname}${window.location.search}`;
  if (currentUrl !== nextUrl || window.location.hash) {
    history.replaceState(null, '', nextUrl);
  }
}

/**
 * @param {{ year: number, month: number, selectedDay: string }} state
 * @returns {boolean}
 */
function applyStateFromLocation(state) {
  const fromQuery = parseViewQuery(new URLSearchParams(window.location.search));
  if (!fromQuery) return false;
  state.year = fromQuery.year;
  state.month = fromQuery.month;
  state.selectedDay = fromQuery.selectedDay;
  return true;
}

/**
 * @param {{ root?: HTMLElement, dataUrl?: string, cacheBust?: number }} options
 */
export async function mountPublishCalendar(options = {}) {
  const root = options.root || document.getElementById('calendar-app');
  if (!root) throw new Error('Missing calendar root element');

  root.innerHTML = '<p class="calendar-loading">Loading calendar…</p>';

  const cacheBust = options.cacheBust ?? Date.now();
  const dataUrl = options.dataUrl || '/calendar/publish-calendar.json';
  const separator = dataUrl.includes('?') ? '&' : '?';
  const requestUrl = `${dataUrl}${separator}v=${cacheBust}`;

  /** @type {CalendarData} */
  let data;
  try {
    const response = await fetch(requestUrl, { cache: 'no-store' });
    if (!response.ok) {
      throw new Error(`Could not load calendar data (${response.status}). Run: make calendar`);
    }
    data = await response.json();
  } catch (error) {
    root.innerHTML = `<div class="calendar-error">${escapeHtml(error instanceof Error ? error.message : String(error))}</div>`;
    return;
  }

  const state = readPersistedCalendarState();

  root.innerHTML = '';
  const shell = document.createElement('div');
  shell.className = 'calendar-shell';
  root.appendChild(shell);

  const header = document.createElement('header');
  header.className = 'calendar-header';
  shell.appendChild(header);

  const titleBlock = document.createElement('div');
  const h1 = document.createElement('h1');
  h1.textContent = 'Publish calendar';
  titleBlock.appendChild(h1);
  const meta = document.createElement('p');
  meta.className = 'calendar-meta';
  meta.textContent = `${data.entries.length} bundles · updated ${formatGeneratedAt(data.generatedAt)}`;
  titleBlock.appendChild(meta);
  header.appendChild(titleBlock);

  const nav = document.createElement('div');
  nav.className = 'calendar-nav';
  const prevBtn = document.createElement('button');
  prevBtn.type = 'button';
  prevBtn.textContent = '←';
  prevBtn.setAttribute('aria-label', 'Previous month');
  const monthLabel = document.createElement('span');
  monthLabel.className = 'calendar-month-label';
  const nextBtn = document.createElement('button');
  nextBtn.type = 'button';
  nextBtn.textContent = '→';
  nextBtn.setAttribute('aria-label', 'Next month');
  const todayBtn = document.createElement('button');
  todayBtn.type = 'button';
  todayBtn.textContent = 'Today';
  nav.append(prevBtn, monthLabel, nextBtn, todayBtn);
  header.appendChild(nav);

  const legend = document.createElement('div');
  legend.className = 'calendar-legend';
  legend.innerHTML = `
    <span class="calendar-legend-item"><span class="calendar-dot calendar-dot--present"></span> Has file</span>
    <span class="calendar-legend-item"><span class="calendar-dot calendar-dot--missing"></span> Missing</span>
  `;
  const channelLegend = document.createElement('span');
  channelLegend.className = 'calendar-legend-item calendar-legend-item--channels';
  channelLegend.textContent = CHANNEL_BOX_ORDER.map(channelLegendLabel).join(' · ');
  channelLegend.title = 'Channel boxes on each post';
  legend.appendChild(channelLegend);
  shell.appendChild(legend);

  const layout = document.createElement('div');
  layout.className = 'calendar-layout';
  shell.appendChild(layout);

  const gridHost = document.createElement('div');
  gridHost.className = 'calendar-grid-host';
  layout.appendChild(gridHost);

  const detailHost = document.createElement('aside');
  detailHost.className = 'calendar-detail';
  layout.appendChild(detailHost);

  function entriesForDay(dayKey) {
    return data.entries.filter((entry) => entry.planned === dayKey);
  }

  function renderGrid() {
    monthLabel.textContent = `${MONTHS[state.month]} ${state.year}`;

    const grid = document.createElement('div');
    grid.className = 'calendar-grid';

    const weekdays = document.createElement('div');
    weekdays.className = 'calendar-weekdays';
    for (const name of WEEKDAYS) {
      const cell = document.createElement('div');
      cell.className = 'calendar-weekday';
      cell.textContent = name;
      weekdays.appendChild(cell);
    }
    grid.appendChild(weekdays);

    const days = document.createElement('div');
    days.className = 'calendar-days';

    const first = new Date(state.year, state.month, 1);
    const startOffset = (first.getDay() + 6) % 7;
    const daysInMonth = new Date(state.year, state.month + 1, 0).getDate();
    const todayKey = formatDayKey(new Date());

    const totalCells = Math.ceil((startOffset + daysInMonth) / 7) * 7;
    for (let i = 0; i < totalCells; i += 1) {
      const dayNum = i - startOffset + 1;
      const cell = document.createElement('button');
      cell.type = 'button';
      cell.className = 'calendar-day';

      let dayKey = '';
      let displayDay = dayNum;
      let inMonth = dayNum >= 1 && dayNum <= daysInMonth;
      if (!inMonth) {
        cell.classList.add('is-outside');
        const adj = dayNum < 1
          ? new Date(state.year, state.month, dayNum)
          : new Date(state.year, state.month + 1, dayNum - daysInMonth);
        dayKey = formatDayKey(adj);
        displayDay = adj.getDate();
      } else {
        dayKey = formatDayKey(new Date(state.year, state.month, dayNum));
      }

      if (dayKey === todayKey) cell.classList.add('is-today');
      if (dayKey === state.selectedDay) cell.classList.add('is-selected');

      const head = document.createElement('div');
      head.className = 'calendar-day-head';

      const weekday = document.createElement('span');
      weekday.className = 'calendar-day-weekday';
      weekday.textContent = weekdayLabel(dayKey);
      head.appendChild(weekday);

      const num = document.createElement('span');
      num.className = 'calendar-day-num';
      num.textContent = String(displayDay);
      head.appendChild(num);

      cell.appendChild(head);

      const items = entriesForDay(dayKey);
      if (items.length > 0) {
        const list = document.createElement('div');
        list.className = 'calendar-day-items';
        const max = 3;
        for (const entry of items.slice(0, max)) {
          const chip = document.createElement('div');
          chip.className = 'calendar-day-chip';

          const chipTitle = document.createElement('span');
          chipTitle.className = 'calendar-day-chip-title';
          chipTitle.textContent = entry.title;
          chip.appendChild(chipTitle);

          chip.appendChild(renderChannelBoxes(entry));

          chip.title = `${entry.bundle}\n${channelSummary(entry)}`;
          list.appendChild(chip);
        }
        if (items.length > max) {
          const more = document.createElement('span');
          more.className = 'calendar-day-more';
          more.textContent = `+${items.length - max} more`;
          list.appendChild(more);
        }
        cell.appendChild(list);
      }

      cell.addEventListener('click', () => {
        state.selectedDay = dayKey;
        render();
      });

      days.appendChild(cell);
    }

    grid.appendChild(days);
    gridHost.innerHTML = '';
    gridHost.appendChild(grid);
  }

  function renderDetail() {
    detailHost.innerHTML = '';
    const entries = entriesForDay(state.selectedDay);

    const h2 = document.createElement('h2');
    h2.textContent = entries.length === 1 ? '1 post' : `${entries.length} posts`;
    detailHost.appendChild(h2);

    const dateLine = document.createElement('p');
    dateLine.className = 'calendar-detail-date';
    dateLine.textContent = formatDisplayDate(state.selectedDay);
    detailHost.appendChild(dateLine);

    if (entries.length === 0) {
      const empty = document.createElement('p');
      empty.className = 'calendar-detail-empty';
      empty.textContent = 'Nothing planned for this day.';
      detailHost.appendChild(empty);
      return;
    }

    const summary = document.createElement('p');
    summary.className = 'calendar-detail-summary';
    summary.textContent = dayReadinessSummary(entries);
    detailHost.appendChild(summary);

    for (const entry of entries) {
      detailHost.appendChild(renderEntry(entry));
    }
  }

  function render() {
    renderGrid();
    renderDetail();
    persistCalendarState(state);
  }

  prevBtn.addEventListener('click', () => {
    state.month -= 1;
    if (state.month < 0) {
      state.month = 11;
      state.year -= 1;
    }
    render();
  });

  nextBtn.addEventListener('click', () => {
    state.month += 1;
    if (state.month > 11) {
      state.month = 0;
      state.year += 1;
    }
    render();
  });

  todayBtn.addEventListener('click', () => {
    const now = new Date();
    state.year = now.getFullYear();
    state.month = now.getMonth();
    state.selectedDay = formatDayKey(now);
    render();
  });

  window.addEventListener('popstate', () => {
    if (applyStateFromLocation(state)) {
      renderGrid();
      renderDetail();
    }
  });

  render();
}

/** @param {CalendarEntry} entry */
function renderEntry(entry) {
  const wrap = document.createElement('article');
  wrap.className = 'calendar-entry';

  const title = document.createElement('h3');
  title.className = 'calendar-entry-title';
  title.textContent = entry.title;
  wrap.appendChild(title);

  const meta = document.createElement('p');
  meta.className = 'calendar-entry-meta';
  const typePart = entry.type ? ` · ${entry.type}` : '';
  meta.textContent = `${entry.section}${typePart}\n${entry.bundle}`;
  wrap.appendChild(meta);

  const boxes = renderChannelBoxes(entry);
  boxes.classList.add('calendar-channel-boxes--detail');
  wrap.appendChild(boxes);

  const channelNotes = document.createElement('p');
  channelNotes.className = 'calendar-entry-channel-notes';
  channelNotes.textContent = channelStatusNotes(entry);
  wrap.appendChild(channelNotes);

  return wrap;
}

/**
 * @param {CalendarEntry[]} entries
 */
function dayReadinessSummary(entries) {
  let present = 0;
  let missing = 0;

  for (const entry of entries) {
    for (const ch of Object.values(entry.channels)) {
      if (ch.status === 'present') present += 1;
      else if (ch.status === 'missing') missing += 1;
    }
  }

  const parts = [];
  if (present) parts.push(`${present} has file`);
  if (missing) parts.push(`${missing} missing`);
  return parts.length > 0 ? parts.join(' · ') : 'No channel status';
}

/** @param {CalendarEntry} entry */
function renderChannelBoxes(entry) {
  const row = document.createElement('div');
  row.className = 'calendar-channel-boxes';
  row.setAttribute('aria-label', channelSummary(entry));

  for (const id of CHANNEL_BOX_ORDER) {
    const ch = entry.channels[id];
    if (!ch) continue;
    const box = document.createElement('span');
    box.className = `calendar-channel-box calendar-channel-box--${ch.status}`;
    box.textContent = CHANNEL_SHORT[id] || id.toUpperCase();
    box.title = `${channelLegendLabel(id)}: ${humanStatus(ch.status)}`;
    row.appendChild(box);
  }

  return row;
}

/** @param {CalendarEntry} entry */
function channelStatusNotes(entry) {
  const parts = [];
  for (const id of CHANNEL_BOX_ORDER) {
    const ch = entry.channels[id];
    if (!ch) continue;
    parts.push(`${channelLegendLabel(id)} ${humanStatus(ch.status)}`);
  }
  return parts.join(' · ');
}

/** @param {CalendarEntry} entry */
function channelSummary(entry) {
  /** @type {Record<string, number>} */
  const counts = {};
  for (const ch of Object.values(entry.channels)) {
    counts[ch.status] = (counts[ch.status] || 0) + 1;
  }
  const parts = [];
  if (counts.present) parts.push(`${counts.present} has file`);
  if (counts.missing) parts.push(`${counts.missing} missing`);
  return parts.length > 0 ? parts.join(' · ') : 'No channel data';
}

/** @param {string} dayKey */
function weekdayLabel(dayKey) {
  const [y, m, d] = dayKey.split('-').map(Number);
  const dayIndex = (new Date(y, m - 1, d).getDay() + 6) % 7;
  return WEEKDAYS[dayIndex];
}

/** @param {Date} date */
function formatDayKey(date) {
  const y = date.getFullYear();
  const m = String(date.getMonth() + 1).padStart(2, '0');
  const d = String(date.getDate()).padStart(2, '0');
  return `${y}-${m}-${d}`;
}

/** @param {string} dayKey */
function formatDisplayDate(dayKey) {
  const [y, m, d] = dayKey.split('-').map(Number);
  const date = new Date(y, m - 1, d);
  return date.toLocaleDateString(undefined, {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

/** @param {string} iso */
function formatGeneratedAt(iso) {
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return iso;
  return date.toLocaleString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

/** @param {string} status */
function humanStatus(status) {
  switch (status) {
    case 'present': return 'has file';
    case 'missing': return 'missing';
    default: return status;
  }
}

/** @param {string} text */
function escapeHtml(text) {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}
