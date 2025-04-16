app.component('Dashboard', {
    template: /*html*/ `
        <div class="container-fluid">
            <div class="row">
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Records">
                        <template #body>
                            <div class="stats-card text-center">
                                <h2>{{ totalRecords }}</h2>
                                <p>Total Records</p>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Avg Response Time">
                        <template #body>
                            <div class="stats-card text-center">
                                <h2>{{ avgResponseTime }}ms</h2>
                                <p>Client-side</p>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Avg Server Time">
                        <template #body>
                            <div class="stats-card text-center">
                                <h2>{{ avgServerTime }}ms</h2>
                                <p>Server-side</p>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Success Rate">
                        <template #body>
                            <div class="stats-card text-center">
                                <h2>{{ successRate }}%</h2>
                                <div class="progress mt-2">
                                    <div class="progress-bar" 
                                         role="progressbar" 
                                         :style="{ width: successRate + '%' }" 
                                         :class="getSuccessRateClass()">
                                    </div>
                                </div>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><session-details-card class="w-100" :session-details="sessionDetails"></session-details-card></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><response-time-chart :records="records"></response-time-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><status-chart :records="records"></status-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><content-types-chart :records="records"></content-types-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><methods-chart :records="records"></methods-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><encoding-types-chart :records="records"></encoding-types-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0">
                    <box maximizable="false" label="Top Requested Paths">
                        <template #body>
                            <div class="stats-table" style="max-height: 250px overflow-y: auto">
                                <table class="table table-sm table-hover m-0">
                                    <thead>
                                        <tr>
                                            <th>Path</th>
                                            <th class="text-end">Count</th>
                                            <th class="text-end">Avg Time</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        <tr v-for="path in topPaths" :key="path.path">
                                            <td class="text-truncate" style="max-width: 200px" :title="path.path">{{ path.path }}</td>
                                            <td class="text-end">{{ path.count }}</td>
                                            <td class="text-end">{{ path.avgTime }}ms</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </template>
                    </box>
                </div>
            </div>
        </div>
    `,
    inject: ['$timer'],
    props: {
        sessionDetails: Object,
        records: Array
    },
    data() {
        return {
            totalRecords: 0,
            avgResponseTime: 0,
            avgServerTime: 0,
            successRate: 0,
            topPaths: [],
            postponeUpdate: this.$timer.debounce(() => this.processRecordsData(), 100)
        }
    },
    mounted() { this.processRecordsData() },
    watch: { records: { handler() { this.postponeUpdate() }, deep: true } },
    methods: {
        processRecordsData() {
            this.totalRecords = this.records.length

            if (this.totalRecords === 0) return

            let totalResponseTime = 0
            let totalServerTime = 0
            let successCount = 0
            let pathMap = {}

            this.records.forEach(record => {
                // Response times
                const responseTime = record.response.elapsed || 0
                totalResponseTime += responseTime

                // Server times
                const serverTime = record.response.serverElapsed || 0
                totalServerTime += serverTime

                // Success rate (codes 2xx)
                if (Math.floor(record.response.code / 100) === 2) {
                    successCount++
                }

                // Path tracking
                const path = record.request.path
                if (!pathMap[path]) {
                    pathMap[path] = { path, count: 0, totalTime: 0 }
                }

                pathMap[path].count++
                pathMap[path].totalTime += responseTime
            })

            this.avgResponseTime = Math.round(totalResponseTime / this.totalRecords)
            this.avgServerTime = Math.round(totalServerTime / this.totalRecords)
            this.successRate = Math.round((successCount / this.totalRecords) * 100)

            // Process top paths
            this.topPaths = Object.values(pathMap)
                .map(item => ({
                    path: item.path,
                    count: item.count,
                    avgTime: Math.round(item.totalTime / item.count)
                }))
                .sort((a, b) => b.count - a.count)
                .slice(0, 10)
        },
        getSuccessRateClass() {
            if (this.successRate >= 90) return 'bg-success'
            if (this.successRate >= 70) return 'bg-warning'
            return 'bg-danger'
        },
    }
})

app.component('SessionDetailsCard', {
    template: /*html*/ `
        <box label="Session Details" maximizable="false">
            <template #body>
                <div class="stats-table" style="max-height: 250px overflow-y: auto">
                    <table class="table table-hover m-0" aria-label="Session Details">
                        <tbody>
                            <tr v-for="(value, key) in sessionDetails">
                                <template v-if='value'>
                                    <td class="highlight">{{ getKey(key) }}</td>
                                    <td v-html="getValue(value)"></td>
                                </template>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </template>
        </box>
    `,
    props: { sessionDetails: Object },
    data() {
        return {
            keyMap: {
                clientId: 'Client ID',
                http: 'URL',
                https: 'Secure URL',
                proxying: 'Proxying',
                sniffer: 'Sniffer',
                tunnel: 'Tunnel URL',
                maxRecords: 'Max Stored Records',
                expireAt: 'Expire At'
            }
        }
    },
    methods: {
        getKey(key) {
            return this.keyMap[key] || key
        },

        getValue(value) {
            if ((value + "").indexOf("://") !== -1) {
                return '<a target="_blank" href="' + value + '">' + value + '</a>'
            }

            return value
        }
    }
})

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
    inject: ['$chart', '$timer'],
    props: { records: Array },
    data() { return { chart: null, postponeUpdate: this.$timer.debounce(() => this.updateChart(), 100) } },
    mounted() { this.createChart() },
    beforeUnmount() { this.chart && this.chart.destroy() },
    watch: { records: { handler() { this.postponeUpdate() }, deep: true } },
    methods: {
        createChart() {
            const data = this.getData()
            const labels = Object.keys(data)
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
                labels,
                Object.values(data),
                backgroundColor,
                {
                    indexAxis: 'y',
                    scales: {
                        x: {
                            beginAtZero: true
                        }
                    }
                }
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
                '101-200ms': 0,
                '501-1000ms': 0,
                '1-2s': 0,
                '2-5s': 0,
                '5s+': 0
            }

            this.records.forEach(record => {
                const time = record.response.elapsed || 0
                if (time <= 100) data['0-100ms']++
                else if (time <= 200) data['101-200ms']++
                else if (time <= 1000) data['501-1000ms']++
                else if (time <= 2000) data['1-2s']++
                else if (time <= 2000) data['2-5s']++
                else data['5s+']++
            })

            return data
        }
    }
})

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
    inject: ['$chart', '$timer'],
    props: { records: Array },
    data() { return { chart: null, postponeUpdate: this.$timer.debounce(() => this.updateChart(), 100) } },
    mounted() { this.createChart() },
    beforeUnmount() { this.chart && this.chart.destroy() },
    watch: { records: { handler() { this.postponeUpdate() }, deep: true } },
    methods: {
        createChart() {
            const data = this.getData()
            const labels = Object.keys(data)
            const backgroundColor = [
                'rgba(200, 200, 200, 0.7)', // 1xx
                'rgba(83, 180, 119, 0.7)',  // 2xx
                'rgba(25, 118, 210, 0.7)',  // 3xx
                'rgba(255, 206, 86, 0.7)',  // 4xx
                'rgba(255, 99, 132, 0.7)'   // 5xx
            ]

            this.chart = this.$chart.create(
                this.$refs.canvas,
                'bar',
                labels,
                Object.values(data),
                backgroundColor,
                {
                    plugins: {
                        tooltip: {
                            callbacks: {
                                footer: (items) => {
                                    const labels = {
                                        '1xx': 'Informational',
                                        '2xx': 'Success',
                                        '3xx': 'Redirection',
                                        '4xx': 'Client Error',
                                        '5xx': 'Server Error'
                                    };
                                    return labels[items[0].label] || '';
                                }
                            }
                        }
                    }
                }
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

app.component('ContentTypesChart', {
    template: /*html*/ `
        <box maximizable="false" label="Content Types">
            <template #body>
                <div class="stats-chart">
                    <canvas ref="canvas"></canvas>
                </div>
            </template>
        </box>
    `,
    inject: ['$chart', '$timer'],
    props: { records: Array },
    data() { return { chart: null, postponeUpdate: this.$timer.debounce(() => this.updateChart(), 100) } },
    mounted() { this.createChart() },
    beforeUnmount() { this.chart && this.chart.destroy() },
    watch: { records: { handler() { this.postponeUpdate() }, deep: true } },
    methods: {
        createChart() {
            const data = this.getData()
            const labels = Object.keys(data)
            const backgroundColor = [
                'rgba(83, 180, 119, 0.7)',
                'rgba(25, 118, 210, 0.7)',
                'rgba(255, 206, 86, 0.7)',
                'rgba(255, 99, 132, 0.7)',
                'rgba(153, 102, 255, 0.7)',
                'rgba(255, 159, 64, 0.7)',
                'rgba(199, 199, 199, 0.7)',
                'rgba(84, 199, 219, 0.7)',
                'rgba(71, 71, 71, 0.7)',
                'rgba(199, 134, 207, 0.7)'
            ]

            this.chart = this.$chart.create(
                this.$refs.canvas,
                'pie',
                labels,
                Object.values(data),
                backgroundColor,
                { plugins: { legend: { position: 'right' } } }
            )
        },

        updateChart() {
            const data = this.getData()
            this.chart.data.labels = Object.keys(data)
            this.chart.data.datasets[0].data = Object.values(data)
            this.chart.update()
        },

        getData() {
            const data = {}
            this.records.forEach(record => {
                const contentType = record.response.header['Content-Type'] || record.response.header['content-type'] || 'unknown'
                data[contentType] = (data[contentType] || 0) + 1
            })
            return data
        }
    }
})

app.component('MethodsChart', {
    template: /*html*/ `
        <box maximizable="false" label="HTTP Methods">
            <template #body>
                <div class="stats-chart">
                    <canvas ref="canvas"></canvas>
                </div>
            </template>
        </box>
    `,
    inject: ['$chart', '$timer'],
    props: { records: Array },
    data() { return { chart: null, postponeUpdate: this.$timer.debounce(() => this.updateChart(), 100) } },
    mounted() { this.createChart() },
    beforeUnmount() { this.chart && this.chart.destroy() },
    watch: { records: { handler() { this.postponeUpdate() }, deep: true  } },
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
                { plugins: { legend: { position: 'right' } } }
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

app.component('EncodingTypesChart', {
    template: /*html*/ `
        <box maximizable="false" label="Encoding Types">
            <template #body>
                <div class="stats-chart">
                    <canvas ref="canvas"></canvas>
                </div>
            </template>
        </box>
    `,
    inject: ['$chart', '$timer'],
    props: { records: Array },
    data() { return { chart: null, postponeUpdate: this.$timer.debounce(() => this.updateChart(), 100) } },
    mounted() { this.createChart() },
    beforeUnmount() { this.chart && this.chart.destroy() },
    watch: { records: { handler() { this.postponeUpdate() }, deep: true  } },
    methods: {
        createChart() {
            const data = this.getData()

            this.chart = this.$chart.create(
                this.$refs.canvas,
                'doughnut',
                Object.keys(data),
                Object.values(data),
                [
                    'rgba(83, 180, 119, 0.7)',
                    'rgba(54, 162, 235, 0.7)',
                    'rgba(255, 206, 86, 0.7)',
                    'rgba(153, 102, 255, 0.7)',
                    'rgba(255, 99, 132, 0.7)'
                ],
                { plugins: { legend: { position: 'right' } } }
            )
        },

        updateChart() {
            const data = this.getData()

            this.chart.data.labels = Object.keys(data)
            this.chart.data.datasets[0].data = Object.values(data)
            
            this.chart.update()
        },

        getData() {
            const data = {}

            this.records.forEach(record => {
                let encoding = 'uncompressed'
                if (record.response.header && (record.response.header['Content-Encoding'] || record.response.header['content-encoding'])) {
                    encoding = record.response.header['Content-Encoding'] || record.response.header['content-encoding']
                }
                data[encoding] = (data[encoding] || 0) + 1
            })

            return data
        }
    }
})