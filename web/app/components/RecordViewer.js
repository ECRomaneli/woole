app.component('RecordViewer', {
    template: /*html*/ `
        <box label-img="request" label="Request" class="w-100">
            <template #buttons>
                <button type="button" class="btn btn-sm" @click="openCurlViewer()">cURL</button>
                <div class="btn-group ms-2 me-2">
                    <button type="button" class="btn btn-sm" @click="$bus.trigger('record.replay', record)">Replay</button>
                    <button type="button" class="btn btn-sm" @click="$refs.reqEditor.show()">w/ Changes</button>
                </div>
            </template>
            <template #body>
                <record-item 
                    :titleGroup="[record.request.method, record.request.url, record.request.proto]" 
                    :item="record.request">
                </record-item>
            </template>
        </box>
        <box label-img="response" label="Response" class="w-100">
            <template #body>
                <record-item 
                    :titleGroup="[record.response.code, $constants.HTTP_STATUS_MESSAGE[record.response.code], record.response.proto]" 
                    :item="record.response">
                </record-item>
            </template>
        </box>

        <request-editor ref="reqEditor" :originalRequest="record.request"></request-editor>
        <code-viewer ref='curlViewer' type="shellscript" :code="record.request.curl"></code-viewer>
    `,
    inject: ['$constants', '$woole'],
    props: { record: Object },

    methods: {
        openCurlViewer() {
            if (!this.record.request.curl) {
                this.$woole.parseRequestToCurl(this.record.request)
            }
            this.$refs.curlViewer.show()
        }
    }
})

app.component('RecordItem', {
    template: /*html*/ `
        <div class="highlighted-group input-group mb-3">
            <span class="input-group-text">{{ titleGroup[0] }}</span>
            <input type="text" class="form-control" disabled :value="titleGroup[1]">
            <span v-if="titleGroup[2]" class="input-group-text">{{ titleGroup[2] }}</span>
        </div>

        <ul class="inline-tabs">
            <li v-if="hasHeader()" @click="tab = 'header'">
                <button class="tab" :class="{ active: tab === 'header' }">Header</button>
            </li>
            <li v-if="hasParam()" @click="tab = 'param'">
                <button class="tab" :class="{ active: tab === 'param' }">Params</button>
            </li>
            <li v-if="hasBody()" @click="tab = 'body'; $refs.codeEditor.forceUpdate()">
                <button class="tab" :class="{ active: tab === 'body' }">Body</button>
            </li>

            <li v-if="isPreviewSupported()" @click="tab = 'preview'">
                <button class="tab" :class="{ active: tab === 'preview' }">Preview</button>
            </li>
        </ul>

        <div class="tab-content">
            <div class="tab-pane show" :class="{ active: tab === 'header' }">
                <map-table :map="item.header"></map-table>
            </div>

            <div class="tab-pane" :class="{ active: tab === 'param' }">
                <map-table :map="item.queryParams"></map-table>
            </div>

            <div class="tab-pane mt-3" :class="{ active: tab === 'body' }">
                <code-editor ref="codeEditor" :type="content.type" :code="item.body" :readOnly="true" :minLines="2" :maxLines="39"></code-editor>
            </div>

            <div class="tab-pane mt-3" :class="{ active: tab === 'preview' }">
                <base64-viewer :category="content.category" :type="content.type" :data="item.b64Body"></base64-viewer>
            </div>
        </div>
    `,
    inject: ['$woole'],
    props: { titleGroup: Array, item: Object },
    data() { return { supportedPreviews: ['image', 'video', 'audio'], tab: 'header' } },
    beforeMount() { this.parseBody() },
    beforeUpdate() {
        this.parseBody()
        this.closeUnavailableTab()
    },
    methods: { 
        closeUnavailableTab() {
            let tab = this.tab
            if (tab === 'preview' && !this.isPreviewSupported()) { tab = 'body' }
            if (tab === 'body' && !this.hasBody()) { tab = 'param' }
            if (tab === 'param' && !this.hasParam()) { tab = 'header' }
            if (tab === 'header' && !this.hasHeader()) { tab = '' }
            this.tab = tab
        },

        parseBody() {
            this.content = {}
            if (this.item.body === void 0 || this.item.body === null) { this.item.body = '' }
            if (!this.item.header) { return }
            const contentType = this.item.header['Content-Type']
            if (contentType) {
                this.content = this.$woole.parseContentType(contentType)
            }
        },

        hasHeader() { return this.item.header && Object.keys(this.item.header).length > 0 },
        hasParam() { return this.item.queryParams && Object.keys(this.item.queryParams).length > 0 },
        hasBody() { return this.item.body !== '' },
        isPreviewSupported() { return this.supportedPreviews.some(c => c === this.content.category) }
    }
    
})