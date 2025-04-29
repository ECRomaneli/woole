app.component('MethodsChart', {
    template: /*html*/ `
        <box maximizable="false" label="HTTP Methods">
            <template #body>
                <div v-show="records.length" class="stats-chart">
                    <canvas ref="canvas"></canvas>
                </div>
                <div v-if="!records.length" class="stats-chart d-flex align-items-center justify-content-center">
                    <span class="h4">NO DATA</span>
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
            const backgroundColor = this.getBackgroundColor(labels)

            this.chart = this.$chart.create(
                this.$refs.canvas,
                'doughnut',
                labels,
                Object.values(data),
                backgroundColor,
                null,
                label => this.$bus.trigger('sidebar.search', `request.method: ${label}`)
            )
        },

        updateChart() {
            const data = this.getData()

            this.chart.data.labels = Object.keys(data)
            this.chart.data.datasets[0].data = Object.values(data)
            this.chart.data.datasets[0].backgroundColor = this.getBackgroundColor(Object.keys(data))

            this.chart.update()
        },

        getBackgroundColor(labels) {
            const methodColors = {
                'GET': 'rgba(83, 180, 119, 0.7)',
                'POST': 'rgba(25, 118, 210, 0.7)',
                'PUT': 'rgba(255, 206, 86, 0.7)',
                'DELETE': 'rgba(255, 99, 132, 0.7)',
                'PATCH': 'rgba(153, 102, 255, 0.7)',
                'HEAD': 'rgba(255, 159, 64, 0.7)',
                'OPTIONS': 'rgba(199, 199, 199, 0.7)'
            }

            return labels.map(method =>
                methodColors[method] || `rgba(${Math.floor(Math.random() * 200)}, ${Math.floor(Math.random() * 200)}, ${Math.floor(Math.random() * 200)}, 0.7)`
            )
        },

        getData() {
            const data = {}

            this.records.forEach(record => {
                const method = record.request.method
                data[method] = (data[method] || 0) + 1
            })

            return data
        }
    }
})
