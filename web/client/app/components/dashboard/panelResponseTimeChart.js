app.component('ResponseTimeChart', {
    template: /*html*/ `
        <box maximizable="false" label="Response Time">
            <template #body>
                <div class="stats-chart">
                    <canvas ref="canvas"></canvas>
                </div>
            </template>
        </box>
    `,
    inject: ['$chart'],
    props: { records: Array },
    data() { return { chart: null, labelToRange: {
        '0-100ms': '0-100ms',
        '101-500ms': '101-500ms',
        '501-1000ms': '501-1000ms',
        '1-2s': '1000-2000ms',
        '2-5s': '2000-5000ms',
        '5s+': '5000ms-'
    } } },
    mounted() { this.createChart() },
    beforeUnmount() { this.chart && this.chart.destroy() },
    watch: { records() { this.updateChart() } },
    methods: {
        createChart() {
            const data = this.getData()
            const backgroundColor = [
                'rgba(83, 180, 119, 0.7)',
                'rgba(105, 192, 150, 0.7)',
                'rgba(255, 206, 86, 0.7)',
                'rgba(255, 159, 64, 0.7)',
                'rgba(255, 129, 102, 0.7)',
                'rgba(255, 99, 132, 0.7)'
            ]

            this.chart = this.$chart.create(
                this.$refs.canvas,
                'bar',
                Object.keys(data),
                Object.values(data),
                backgroundColor,
                { indexAxis: 'y', scales: { x: { beginAtZero: true } }, plugins: { legend: { display: false } } },
                label => this.$bus.trigger('sidebar.search', `response.elapsed~: ${this.labelToRange[label]}`)
            )
        },

        updateChart() {
            const data = this.getData()
            this.chart.data.labels = Object.keys(data)
            this.chart.data.datasets[0].data = Object.values(data)
            this.chart.update()
        },

        getData() {
            const data = {
                '0-100ms': 0,
                '101-500ms': 0,
                '501-1000ms': 0,
                '1-2s': 0,
                '2-5s': 0,
                '5s+': 0
            }

            this.records.forEach(record => {
                const time = record.response.elapsed || 0
                if (time <= 100) data['0-100ms']++
                else if (time <= 500) data['101-500ms']++
                else if (time <= 1000) data['501-1000ms']++
                else if (time <= 2000) data['1-2s']++
                else if (time <= 5000) data['2-5s']++
                else data['5s+']++
            })

            return data
        }
    }
})
