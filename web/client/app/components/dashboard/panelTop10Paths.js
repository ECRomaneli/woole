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
                            <tr v-for="path in data" :key="path.path" @click="searchPath(path.path)">
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
    inject: ['$woole'],
    props: { records: Array },
    data() { return { data: {} } },
    mounted() { this.getData() },
    watch: { records() { this.getData() } },
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
        },

        searchPath(path) {
            this.$bus.trigger('sidebar.search', 'request.path *: "^' + this.$woole.escapeRegex(path) + '$"')
        }
    }
})
