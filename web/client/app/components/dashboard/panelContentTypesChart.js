app.component('ContentTypesChart', {
    template: /*html*/ `
        <box maximizable="false" label="Status Codes">
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
                'pie',
                Object.keys(data),
                Object.values(data),
                null,
                { plugins: { legend: { position: 'right' } } },
                label => this.$bus.trigger('sidebar.search', label === 'no content-type' ? 'not response.header.content-type' : `response.header.content-type: ${label}`)
            )

            this.$chart.colorfy(this.chart)
            this.chart.update()
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
                const contentType = record.response.getHeader('Content-Type', 'no content-type').split(';')[0].trim()
                data[contentType] = (data[contentType] || 0) + 1
            })
            return data
        }
    }
})
