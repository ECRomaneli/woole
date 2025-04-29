app.component('RemoteAddressDetails', {
    template: /*html*/ `
        <box maximizable="false" label="Remote Address Details">
            <template #body>
                <div class="stats-table remote-address-table">
                    <table class="table table-striped table-hover">
                        <thead>
                            <tr>
                                <th>Address</th>
                                <th class="w-100">Paths</th>
                                <th>Records</th>
                                <th>Size</th>
                                <th>Avg Time</th>
                                <th>Avg Server Time</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(data, ip) in data" :key="ip"  @click="searchIp(ip)">
                                <td>{{ ip }}</td>
                                <td class="paths-column">{{ data.paths }}</td>
                                <td>{{ data.count }}</td>
                                <td>{{ data.totalSize }}</td>
                                <td>{{ data.avgResponseTime }}ms</td>
                                <td>{{ data.avgServerTime }}ms</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </template>
        </box>
    `,
    inject: ['$woole'],
    props: { records: Array },
    data() {
        return {
            data: {}
        }
    },
    mounted() {
        this.updateData()
    },
    watch: {
        records: {
            handler() {
                this.updateData()
            },
            deep: true
        }
    },
    methods: {
        updateData() {
            this.data = this.getData()
        },

        searchIp(ip) {
            this.$bus.trigger('sidebar.search', `remoteAddr*: "^\\[?${this.$woole.escapeRegex(ip)}(]|:|$)"`)
        },
        
        getData() {
            const ipData = {}

            this.records.forEach(record => {
                if (!record.request.remoteAddr || !record.request.path) return

                const ip = this.$woole.parseAddress(record.request.remoteAddr)?.ip
                const path = record.request.path
                const responseTime = record.response.elapsed || 0
                const serverTime = record.response.serverElapsed || 0
                const contentLength = parseInt(record.response.getHeader('Content-Length', 0), 10)

                data = ipData[ip]

                if (!data) {
                    data = {
                        paths: [],
                        totalResponseTime: 0,
                        totalServerTime: 0,
                        totalSize: 0,
                        count: 0
                    }
                    ipData[ip] = data
                }

                if (data.paths.length < 10 && !data.paths.includes(path)) {
                    data.paths.push(path)
                }

                data.totalResponseTime += responseTime
                data.totalServerTime += serverTime
                data.totalSize += contentLength
                data.count += 1
            })

            const result = {}
            Object.entries(ipData).forEach(([ip, data]) => {
                result[ip] = {
                    paths: data.paths.join(', ').slice(0, 360),
                    count: data.count,
                    totalSize: this.$woole.parseSize(data.totalSize),
                    avgResponseTime: Math.round(data.totalResponseTime / data.count),
                    avgServerTime: Math.round(data.totalServerTime / data.count)
                }
            })

            return result
        }
    }
})