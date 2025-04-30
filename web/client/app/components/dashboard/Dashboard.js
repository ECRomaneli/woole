app.component('Dashboard', {
    template: /*html*/ `
        <div id="dashboard" class="container-fluid pt-2">
            <div class="row">
                 <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Client ID">
                        <template #body>
                            <div class="stats-card">
                                <span class="h5">{{ clientId }}</span>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="URL">
                        <template #body>
                            <div class="stats-card">
                                <a class="h6 m-1" v-if="httpsUrl || httpUrl" :href="httpsUrl || httpUrl" target="_blank">{{ httpsUrl || httpUrl }}</a>
                                <a v-if="httpsUrl" class="h6" :href="httpUrl" target="_blank">{{ httpUrl }}</a>
                                <span class="fw-light" v-else>No HTTPS URL</span>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Tunnel URL">
                        <template #body>
                            <div class="stats-card">
                                <span v-if="tunnelUrl" class="h5">{{ tunnelUrl }}</span>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Expire Date">
                        <template #body>
                            <div class="stats-card">
                                <span class="h5 m-0" v-if="expireDate">{{ expireDate }}</span>
                                <span class="fw-light" v-if="expireRemaining !== null">Expires in {{ expireRemaining | 0 }} minutes</span>
                                <span class="fw-light" v-else-if="expireDate === $constants.NEVER_EXPIRE_MESSAGE">No Expiration</span>
                                <span class="fw-light" v-else>Tunnel is no longer connected</span>
                            </div>
                        </template>
                    </box>
                </div>

                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Records">
                        <template #body>
                            <div class="stats-card">
                                <span class="h4 m-0">{{ totalRecords }} / {{ maxRecords }}</span>
                                <span class="fw-light">Total Records</span>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Avg Response Time">
                        <template #body>
                            <div class="stats-card">
                                <span class="h4 m-0">{{ avgResponseTime }}ms</span>
                                <span class="fw-light">Client-side</span>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Avg Server Time">
                        <template #body>
                            <div class="stats-card">
                                <span class="h4 m-0">{{ avgServerTime }}ms</span>
                                <span class="fw-light">Server-side</span>
                            </div>
                        </template>
                    </box>
                </div>
                <div class="col-xl-3 col-lg-6 p-0">
                    <box maximizable="false" label="Session Status">
                        <template #body>
                            <div :class="'stats-card text-' + (sessionStatus?.color || 'none')">
                                <span class="h5">{{ sessionStatus?.name || '-' }}</span>
                            </div>
                        </template>
                    </box>
                </div>

                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><content-types-chart :records="records"></content-types-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><methods-chart :records="records"></methods-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><encoding-types-chart :records="records"></encoding-types-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><response-time-chart :records="records"></response-time-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><status-chart :records="records"></status-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><top-remote-addrs-chart :records="records"></top-remote-addrs-chart></div>
                <div class="col-md-12 col-lg-6 col-xl-4 p-0"><top-paths-list :records="records"></top-paths-list></div>
                <div class="col-md-12 col-lg-6 col-xl-8 p-0"><remote-address-list :records="records"></remote-address-list></div>
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
            clientId: null,
            httpUrl: null,
            httpsUrl: null,
            tunnelUrl: null,
            maxRecords: null,
            sessionStatus: null,
            expireRemaining: null,
            expireInterval: null,
            expireDate: null,
            delayedRecords: [],
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
        records: { handler() { this.processRecordsData() }, deep: true }
    },
    methods: {
        processRecordsData() {
            this.totalRecords = this.records.length

            let totalResponseTime = 0
            let totalServerTime = 0

            this.records.forEach(record => {
                totalResponseTime += record.response.elapsed || 0
                totalServerTime += record.response.serverElapsed || 0
            })

            this.avgResponseTime = Math.round(totalResponseTime / (this.totalRecords || 1))
            this.avgServerTime = Math.round(totalServerTime / (this.totalRecords || 1))
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

            if (!expireTime) {
                this.expireDate = '-'
                this.expireRemaining = null
                return
            }

            const currentTime = Date.now()

            if (expireTime < currentTime) {
                this.expireDate = 'Expired'
                this.expireRemaining = null
                return
            }

            this.expireDate = new Date(expireTime).toLocaleString()
            this.expireRemaining = Math.max(0, (expireTime - currentTime) / 60000)

            this.expireInterval = setInterval(() => {
                this.expireRemaining = Math.max(0, (expireTime - Date.now()) / 60000)
                if (this.expireRemaining <= 0) { this.setExpireAt() }
            }, 60000)
        }
    }
})
