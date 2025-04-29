app.component('EncodingTypesChart', {
    template: /*html*/ `
        <box maximizable="false" label="Encoding Types">
            <template #body>
                <div class="stats-chart d-flex align-items-center justify-content-center">
                    <canvas v-show="records.length" ref="canvas"></canvas>
                    <span v-if="!records.length" class="h4">NO DATA</span>
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

            this.chart = this.$chart.create(
                this.$refs.canvas,
                'doughnut',
                Object.keys(data),
                Object.values(data),
                null,
                null,
                label => this.$bus.trigger('sidebar.search', label === 'uncompressed' ? 'not content-encoding' : `content-encoding: ${label}`)
            )
            this.$chart.colorfy(this.chart)
        },

        updateChart() {
            const data = this.getData()

            this.chart.data.labels = Object.keys(data)
            this.chart.data.datasets[0].data = Object.values(data)
            this.$chart.colorfy(this.chart)

            this.chart.update()
        },

        getData() {
            const data = {}

            this.records.forEach(record => {
                const encoding = record.response.getHeader('Content-Encoding', 'uncompressed')
                data[encoding] = (data[encoding] || 0) + 1
            })

            return data
        }
    }
})
