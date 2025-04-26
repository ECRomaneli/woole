app.component('StatusChart', {
    template: /*html*/ `
        <box maximizable="false" label="Status Codes">
            <template #body>
                <div class="stats-chart">
                    <canvas ref="canvas"></canvas>
                </div>
            </template>
        </box>
    `,
    inject: ['$chart'],
    props: { records: Array },
    data() { return { chart: null } },
    mounted() { this.createChart() },
    beforeUnmount() { this.chart && this.chart.destroy() },
    watch: { records() { this.updateChart() } },
    methods: {
        createChart() {
            const data = this.getData()
            const labels = Object.keys(data)
            const backgroundColor = [
                'rgba(200, 200, 200, 0.7)', // 1xx
                'rgba(83, 180, 119, 0.7)',  // 2xx
                'rgba(210, 145, 25, 0.7)',  // 3xx
                'rgba(255, 86, 86, 0.7)',   // 4xx
                'rgba(255, 99, 177, 0.7)',  // 5xx
                'rgba(32, 32, 32, 0.7)'     // Unknown
            ]

            this.chart = this.$chart.create(
                this.$refs.canvas,
                'bar',
                labels,
                Object.values(data),
                backgroundColor,
                {
                    plugins: {
                        tooltip: { callbacks: { footer: (items) => {
                            const labels = {
                                '1xx': 'Informational',
                                '2xx': 'Success',
                                '3xx': 'Redirection',
                                '4xx': 'Client Error',
                                '5xx': 'Server Error',
                                '9xx': 'Unknown'
                            }
                            return labels[items[0].label] || ''
                        } } },
                        legend: { display: false }
                    }
                },
                label => this.$bus.trigger('sidebar.search', `response.codeGroup: ${label}`)
            )
        },

        updateChart() {
            const data = this.getData()
            this.chart.data.labels = Object.keys(data)
            this.chart.data.datasets[0].data = Object.values(data)
            this.chart.update()
        },

        getData() {
            const data = { '1xx': 0, '2xx': 0, '3xx': 0, '4xx': 0, '5xx': 0 }
            this.records.forEach(record => {
                const group = `${Math.floor(record.response.code / 100)}xx`
                if (data[group] !== undefined) {
                    data[group]++
                }
            })
            return data
        }
    }
})
