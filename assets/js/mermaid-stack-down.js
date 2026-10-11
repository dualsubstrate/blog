/**
 * Enlarge only vertical stack edges in Mermaid flowcharts (subgraph → subgraph).
 * Horizontal node→note links stay default. Lives in the site repo so the LoveIt submodule stays clean.
 */
(function () {
    function enhanceMermaidStackDownArrows(container) {
        var svg = container.querySelector('svg');
        if (!svg) return;
        var defs = svg.querySelector('defs');
        if (!defs) return;
        var paths = svg.querySelectorAll('g.edgePaths path');
        if (!paths.length) return;

        var largeId = 'mermaid-stack-down-arrow';
        var baseMarker =
            svg.querySelector('marker[id*="flowchart-pointEnd"]') ||
            svg.querySelector('marker[id*="arrowhead"]') ||
            defs.querySelector('marker');
        if (baseMarker && !svg.getElementById(largeId)) {
            var m = baseMarker.cloneNode(true);
            m.id = largeId;
            var scale = 1.85;
            var mw = parseFloat(m.getAttribute('markerWidth') || '12');
            var mh = parseFloat(m.getAttribute('markerHeight') || '12');
            m.setAttribute('markerWidth', String(mw * scale));
            m.setAttribute('markerHeight', String(mh * scale));
            defs.appendChild(m);
        }

        paths.forEach(function (path) {
            var bb;
            try {
                bb = path.getBBox();
            } catch (_e) {
                return;
            }
            if (!(bb.height > bb.width * 1.22 && bb.height > 14)) return;
            path.classList.add('mermaid-stack-down');
            if (svg.getElementById(largeId)) {
                path.setAttribute('marker-end', 'url(#' + largeId + ')');
            }
        });
    }

    function scan() {
        document.querySelectorAll('.mermaid').forEach(function (el) {
            if (!el.querySelector('svg g.edgePaths')) return;
            requestAnimationFrame(function () {
                enhanceMermaidStackDownArrows(el);
            });
        });
    }

    var scheduled = false;
    function scheduleScan() {
        if (scheduled) return;
        scheduled = true;
        requestAnimationFrame(function () {
            scheduled = false;
            scan();
        });
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', function () {
            scan();
            new MutationObserver(scheduleScan).observe(document.body, { childList: true, subtree: true });
        });
    } else {
        scan();
        new MutationObserver(scheduleScan).observe(document.body, { childList: true, subtree: true });
    }
})();
