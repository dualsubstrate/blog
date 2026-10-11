/**
 * When window.config.data[id] is missing, LoveIt theme passes undefined into mermaid.render and Mermaid throws.
 * layouts/partials/plugin/mermaid.html duplicates the source in a hidden textarea; this patch reads it so the
 * theme submodule does not need custom theme.js.
 */
(function () {
    function applyPatch() {
        if (typeof mermaid === 'undefined' || typeof mermaid.render !== 'function' || mermaid.render._siteSourcePatched) {
            return;
        }
        var original = mermaid.render.bind(mermaid);
        mermaid.render = function () {
            var args = Array.prototype.slice.call(arguments);
            var diagramId = args[0];
            var text = args[1];
            if (text === undefined || text === null || String(text).trim() === '') {
                var prefix = 'mermaid-svg-';
                if (typeof diagramId === 'string' && diagramId.indexOf(prefix) === 0) {
                    var domId = diagramId.slice(prefix.length);
                    var ta = document.getElementById(domId + '-mermaid-src');
                    if (ta) {
                        args[1] = (ta.value || ta.textContent || '').trim();
                    }
                }
            }
            return original.apply(mermaid, args);
        };
        mermaid.render._siteSourcePatched = true;
    }

    applyPatch();
})();
