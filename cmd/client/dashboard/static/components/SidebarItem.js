app.vue.component('SidebarItem', {
    template: `
        <div class="list-group-item list-group-item-action py-3 lh-tight">
            <div class="d-flex w-100 align-items-center justify-content-between">
                <div class="mb-1">
                    <span class="badge mr-3" :class="protoBadge()">{{ protocol }}</span>
                    <span class="badge mr-3" :class="methodBadge()">{{ request.method }}</span>
                    <span class="badge mr-3" :class="statusBadge()">{{ response.code }}</span>
                </div>
                <small>{{ record.elapsed }}ms</small>
            </div>
            <div class="mb-1 small">
                <span class="request-path">{{ ellipsis(request.path) }}</span>
                <span v-if="request.query !== void 0" class="request-query badge" :title="request.query">?</span>
            </div>
        </div>
    `,
    props: {
        record: Object
    },
    data() {
        return {
            request: this.record.request,
            response: this.record.response,
            maxLength: 30
        }
    },
    beforeMount() {
        this.request.query = this.request.url.split('?')[1]
    },
    methods: {
        methodBadge() {
            return "bg-" + this.request.method.toLowerCase()
        },
        protoBadge() {
            this.protocol = this.request.proto.split("/")[0]
            return "bg-" + this.protocol.toLowerCase()
        },
        statusBadge() {
            return "bg-status-" + parseInt(this.response.code/100)
        },
        ellipsis(path) {
            let maxLength = this.maxLength
            let hasQuery = this.request.query !== void 0

            if (hasQuery) { maxLength -= 4 }
            let result = path.length < maxLength ? path : '...' + path.substring(path.length - maxLength)
            return result + (hasQuery ? ' ' : '')
        }
    }
})