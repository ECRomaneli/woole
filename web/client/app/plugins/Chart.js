app.use({
    install: (app) => {
        const BACKGROUND_COLORS = [
            'rgba(83, 180, 119, 0.7)',
            'rgba(25, 118, 210, 0.7)',
            'rgba(255, 206, 86, 0.7)',
            'rgba(255, 99, 132, 0.7)',
            'rgba(153, 102, 255, 0.7)',
            'rgba(255, 159, 64, 0.7)',
            'rgba(0, 150, 136, 0.7)',
            'rgba(183, 28, 28, 0.7)',
            'rgba(205, 220, 57, 0.7)',
            'rgba(121, 85, 72, 0.7)',
            'rgba(0, 188, 212, 0.7)',    
            'rgba(216, 27, 96, 0.7)',    
            'rgba(46, 125, 50, 0.7)',
            'rgba(255, 138, 101, 0.7)',  
            'rgba(186, 104, 200, 0.7)',
            'rgba(63, 81, 181, 0.7)',    
            'rgba(128, 203, 196, 0.7)'
        ]

        const DEFAULT_OPTIONS = {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: true,
                    position: 'right',
                    labels: {
                        boxWidth: 12
                    }
                }
            }
        }

        const labelColorByChart = new Map()

        function create(canvas, type, labels, data, colors, options, onClick) {
            options = options || {}
    
            const chart = Vue.markRaw(new Chart(canvas, {
                type: type,
                data: {
                    labels: labels,
                    datasets: [{
                        data: data,
                        backgroundColor: colors || BACKGROUND_COLORS,
                        borderWidth: 0
                    }]
                },
                options: mergeOptions(options)
            }))
            options.onClick = options.onClick || createOnLabelClick(chart, onClick)
    
            return chart
        }

        function colorfy(chart) {
            let labelColors = labelColorByChart.get(chart)

            if (labelColors === void 0) {
                labelColorByChart.set(chart, new Map())
                labelColors = labelColorByChart.get(chart)
            }

            chart.data.datasets[0].backgroundColor = chart.data.labels.map(l => {
                if (!labelColors.has(l)) {
                    labelColors.set(l, BACKGROUND_COLORS[labelColors.size % BACKGROUND_COLORS.length])
                }
                return labelColors.get(l)
            })
        }

        function mergeOptions(userOptions, defaultOptions = DEFAULT_OPTIONS) {
            const mergedOptions = { ...defaultOptions }
        
            for (const key in userOptions) {
                if (userOptions[key] && typeof userOptions[key] === 'object' && !Array.isArray(userOptions[key])) {
                    mergedOptions[key] = mergeOptions(defaultOptions[key] || {}, userOptions[key])
                } else {
                    mergedOptions[key] = userOptions[key]
                }
            }
        
            return mergedOptions;
        }

        function createOnLabelClick(chart, onClick) {
            chart.options.onClick = (_, el) => {
                if (el.length > 0) {
                    const label = chart.data.labels[el[0].index]
                    onClick && onClick(label)
                }
            }
        }

        app.provide('$chart', { create: create, colorfy: colorfy })
    }
})