app.provide('$chart', {
    create(canvas, type, labels, data, backgroundColor, options, onClick) {
        options = options || {}
        options.responsive = options.responsive || true
        options.maintainAspectRatio = options.maintainAspectRatio || false
        options.plugins = options.plugins || {}
        options.plugins.legend = options.plugins.legend || false
        options.onClick = options.onClick || ((_, el) => {
            if (el.length > 0) {
                const label = chart.data.labels[el[0].index]
                onClick && onClick(label)
            }
        })

        const chart = Vue.markRaw(new Chart(canvas, {
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

        return chart
    }
})