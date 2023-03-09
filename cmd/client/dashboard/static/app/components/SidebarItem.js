app.component('SidebarItem', {
    template: /*html*/ `
        <div class="list-group-item list-group-item-action py-3 lh-tight">
            <div class="d-flex w-100 align-items-center justify-content-between">
                <div class="mb-1">
                    <div v-if="isReplay()" class="bg-replay-badge replay-badge badge mr-4"><img src="assets/images/play.svg" alt="replay" /></div>
                    <span class="badge mr-4" :class="methodBadge()">{{ request.method }}</span>
                    <span class="badge" :class="statusBadge()">{{ response.code }}</span>
                </div>
                <div>
                    <small style="font-size: 10px; margin-right: 3px; color: #bbb;">77ms /</small>
                    <small>{{ record.elapsed }}ms</small>
                </div>
                
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
        statusBadge() {
            return "bg-status-" + parseInt(this.response.code/100)
        },
        ellipsis(path) {
            let maxLength = this.maxLength
            let hasQuery = this.request.query !== void 0

            if (hasQuery) { maxLength -= 4 }
            let result = path.length < maxLength ? path : '...' + path.substring(path.length - maxLength)
            return result + (hasQuery ? ' ' : '')
        },
        isReplay() {
            return this.record.id.indexOf('R') !== -1
        }
    }
})