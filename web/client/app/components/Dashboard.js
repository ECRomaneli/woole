app.component('Dashboard', {
    template: /*html*/ `
        <div class="container-fluid">
            <div class="row">
                 <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Client ID">
                        <template #body>
                            <div class="stats-card d-flex align-items-center justify-content-center">
                                <span class="h4">{{ clientId }}</span>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="URL">
                        <template #body>
                            <div class="stats-card text-center">
                                <div class="mb-1">
                                    <a class="h6" v-if="httpsUrl || httpUrl" :href="httpsUrl || httpUrl" target="_blank">{{ httpsUrl || httpUrl }}</a>
                                </div>
                                <a v-if="httpsUrl" class="h6" :href="httpUrl" target="_blank">{{ httpUrl }}</a>
                                <p v-else>No HTTPS URL</p>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Tunnel URL">
                        <template #body>
                            <div class="stats-card d-flex align-items-center justify-content-center">
                                <span v-if="tunnelUrl" class="h4">{{ tunnelUrl }}</span>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Expire Date">
                        <template #body>
                            <div class="stats-card text-center">
                                <span class="h4" v-if="expireDate">{{ expireDate }}</span>
                                <p v-if="expireRemaining !== null">Expires in {{ expireRemaining | 0 }} minutes</p>
                                <p v-else>No Expiration</p>
                            </div>
                        </template>
                    </box>
                </div>

                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Records">
                        <template #body>
                            <div class="stats-card text-center">
                                <span class="h3">{{ totalRecords }} / {{ maxRecords }}</span>
                                <p>Total Records</p>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Avg Response Time">
                        <template #body>
                            <div class="stats-card text-center">
                                <span class="h3">{{ avgResponseTime }}ms</span>
                                <p>Client-side</p>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Avg Server Time">
                        <template #body>
                            <div class="stats-card text-center">
                                <span class="h3">{{ avgServerTime }}ms</span>
                                <p>Server-side</p>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Session Status">
                        <template #body>
                            <div :class="'stats-card d-flex align-items-center text-' + (sessionStatus?.color || 'none') + ' justify-content-center'">
                                <h2>{{ sessionStatus?.name || '-' }}</h2>
                            </div>
                        </template>
                    </box>
                <!-- <div class="col-xl-3 col-lg-6 p-0">
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
                    </box> -->
                </div>

                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><content-types-chart :records="records"></content-types-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><methods-chart :records="records"></methods-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><encoding-types-chart :records="records"></encoding-types-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><response-time-chart :records="records"></response-time-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><status-chart :records="records"></status-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><top-10-paths :records="records"></top-10-paths></div>
            </div>
        </div>
    `,
    inject: ['$timer', '$constants'],
    props: {
        sessionDetails: { type: Object, default: () => ({}) },
        records: Array
    },
    data() {
        return {
            totalRecords: 0,
            avgResponseTime: 0,
            avgServerTime: 0,
            successRate: 0,
            clientId: null,
            httpUrl: null,
            httpsUrl: null,
            tunnelUrl: null,
            maxRecords: null,
            expireDate: null,
            sessionStatus: null,
            expireRemaining: null,
            expireInterval: null,
            postponeUpdate: this.$timer.debounce(() => this.processRecordsData(), 100)
        }
    },
    mounted() { 
        this.processRecordsData()
        this.loadSessionDetails()
    },
    beforeUnmount() {
        if (this.expireInterval) {
            clearInterval(this.expireInterval)
            this.expireInterval = null
        }
    },
    watch: {
        sessionDetails: { handler() { this.loadSessionDetails() }, deep: true },
        records: { handler() { this.postponeUpdate() }, deep: true }
    },
    methods: {
        processRecordsData() {
            this.totalRecords = this.records.length

            let totalResponseTime = 0
            let totalServerTime = 0
            let successCount = 0

            this.records.forEach(record => {
                totalResponseTime += record.response.elapsed || 0
                totalServerTime += record.response.serverElapsed || 0
                if (Math.floor(record.response.code / 100) === 2) { successCount++ }
            })

            this.avgResponseTime = Math.round(totalResponseTime / (this.totalRecords || 1))
            this.avgServerTime = Math.round(totalServerTime / (this.totalRecords || 1))
            this.successRate = Math.round((successCount / (this.totalRecords || 1)) * 100)
        },
        loadSessionDetails() {
            this.clientId = this.sessionDetails.clientId || '-'
            this.httpUrl = this.sessionDetails.http || '-'
            this.httpsUrl = this.sessionDetails.https
            this.tunnelUrl = this.sessionDetails.tunnel || '-'
            this.maxRecords = this.sessionDetails.maxRecords || '∞'
            this.setSessionStatus()
            this.setExpireAt()
        },
        setSessionStatus() {
            this.sessionStatus = { name: this.sessionDetails.status || this.$constants.SESSION_STATUS.CONNECTING }

            switch (this.sessionStatus.name) {
                case this.$constants.SESSION_STATUS.CONNECTING:     this.sessionStatus.color = 'info'; break
                case this.$constants.SESSION_STATUS.CONNECTED:      this.sessionStatus.color = 'success'; break
                case this.$constants.SESSION_STATUS.DISCONNECTED:   this.sessionStatus.color = 'danger'; break
                case this.$constants.SESSION_STATUS.RECONNECTING:   this.sessionStatus.color = 'warning'; break
                default:                                            this.sessionStatus.color = 'none'; break
            }
        },
        setExpireAt() {
            if (this.expireInterval) {
                clearInterval(this.expireInterval)
                this.expireInterval = null
            }

            if (this.sessionDetails.expireAt === this.$constants.NEVER_EXPIRE_MESSAGE) {
                this.expireDate = this.$constants.NEVER_EXPIRE_MESSAGE
                this.expireRemaining = null
                return
            }

            const expireTime = new Date(this.sessionDetails.expireAt).getTime()
            const currentTime = Date.now()

            if (expireTime < currentTime) {
                this.expireDate = 'expired'
                this.expireRemaining = null
            }

            this.expireDate = new Date(expireTime).toLocaleString()
            this.expireRemaining = Math.max(0, (expireTime - currentTime) / 60000)

            this.expireInterval = setInterval(() => {
                this.expireRemaining = Math.max(0, (expireTime - Date.now()) / 60000)
                if (this.expireRemaining <= 0) { this.setExpireAt() }
            }, 60000)
        },
        getSuccessRateClass() {
            if (this.successRate >= 90) return 'bg-success'
            if (this.successRate >= 70) return 'bg-warning'
            return 'bg-danger'
        }
    }
})

app.component('Top10Paths', {
    template: /*html*/ `
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
                            <tr v-for="path in data" :key="path.path" @click="$bus.trigger('sidebar.search', 'request.path *: ^' + path.path + '$')">
                                <td class="text-truncate" style="max-width: 200px" :title="path.path">{{ path.path }}</td>
                                <td class="text-end">{{ path.count }}</td>
                                <td class="text-end">{{ path.avgTime }}ms</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </template>
        </box>
    `,
    props: { records: Array },
    data() { return { data: {} } },
    mounted() { this.getData() },
    watch: { records: { handler() { this.getData() }, deep: true } },
    methods: {
        getData() {
            const pathMap = {}

            this.records.forEach(record => {
                // Path tracking
                const path = record.request.path
                if (!pathMap[path]) {
                    pathMap[path] = { path, count: 0, totalTime: 0 }
                }

                pathMap[path].count++
                pathMap[path].totalTime += record.response.elapsed
            })

            // Process top paths
            this.data = Object.values(pathMap)
                .map(item => ({
                    path: item.path,
                    count: item.count,
                    avgTime: Math.round(item.totalTime / item.count)
                }))
                .sort((a, b) => b.count - a.count)
                .slice(0, 10)
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
    data() { return { chart: null, postponeUpdate: this.$timer.debounce(() => this.updateChart(), 100), labelToRange: {
        '0-100ms': '0-100ms',
        '101-500ms': '101-500ms',
        '501-1000ms': '501-1000ms',
        '1-2s': '1000-2000ms',
        '2-5s': '2000-5000ms',
        '5s+': '5000ms-'
    } } },
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
                { indexAxis: 'y', scales: { x: { beginAtZero: true } } },
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
                    plugins: { tooltip: { callbacks: { footer: (items) => {
                        const labels = {
                            '1xx': 'Informational',
                            '2xx': 'Success',
                            '3xx': 'Redirection',
                            '4xx': 'Client Error',
                            '5xx': 'Server Error',
                            '9xx': 'Unknown'
                        }
                        return labels[items[0].label] || ''
                    } } } }
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
                { plugins: { legend: { position: 'right' } } },
                label => this.$bus.trigger('sidebar.search', `response.header.content-type: ${label}`)
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
                const contentType = record.response.getHeader('Content-Type', 'unknown')
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
                { plugins: { legend: { position: 'right' } } },
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
                { plugins: { legend: { position: 'right' } } },
                label => this.$bus.trigger('sidebar.search', label === 'uncompressed' ? 'not content-encoding' : `content-encoding: ${label}`)
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
                const encoding = record.response.getHeader('Content-Encoding', 'uncompressed')
                data[encoding] = (data[encoding] || 0) + 1
            })

            return data
        }
    }
})