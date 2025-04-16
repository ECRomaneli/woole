app.provide('$chart', {
    create(canvas, type, labels, data, backgroundColor, options) {
        options = options || {}
        options.responsive = options.responsive || true
        options.maintainAspectRatio = options.maintainAspectRatio || false
        options.plugins = options.plugins || {}
        options.plugins.legend = options.plugins.legend || false

        return Vue.markRaw(new Chart(canvas, {
            type: type,
            data: {
                labels: labels,
                datasets: [{
                    data: data,
                    backgroundColor: backgroundColor,
                    borderWidth: 0
                }]
            },
            options: options
        }))
    }
})