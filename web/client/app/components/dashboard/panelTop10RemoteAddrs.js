app.component('Top10RemoteAddrs', {
    template: /*html*/ `
        <box maximizable="false" label="Status Codes">
            <template #body>
                <div class="stats-chart">
                    <canvas ref="canvas"></canvas>
                </div>
            </template>
        </box>
    `,
    inject: ['$chart', '$woole'],
    props: { records: Array },
    data() { return { chart: null } },
    mounted() { this.createChart() },
    beforeUnmount() { this.chart && this.chart.destroy() },
    watch: { records() { this.updateChart() } },
    methods: {
        createChart() {
            const data = this.getData()
            const labels = Object.keys(data)

            this.chart = this.$chart.create(
                this.$refs.canvas,
                'pie',
                labels,
                Object.values(data),
                null,
                null,
                ip => this.$bus.trigger('sidebar.search', `remoteAddr*: "^\\[?${this.$woole.escapeRegex(ip)}(]|:|$)"`)
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
            const ipCounts = {}
            
            this.records.forEach(record => {
                if (!record.request.remoteAddr) { return }

                // Extract IP address without port
                const address = this.$woole.parseAddress(record.request.remoteAddr)
                if (!address?.ip) { return }
                
                ipCounts[address.ip] = (ipCounts[address.ip] || 0) + 1
            })
            
            // Sort by count and take top 10
            return Object.entries(ipCounts)
                .sort((a, b) => b[1] - a[1])
                .slice(0, 10)
                .reduce((obj, [ip, count]) => {
                    obj[ip] = count
                    return obj
                }, {})
        }
    }
})
