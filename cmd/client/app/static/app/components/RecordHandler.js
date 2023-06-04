app.component('RecordViewer', {
    template: /*html*/ `
        <box label-img="request" label="Request" class="me-2 mb-2">
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
        <box label-img="response" label="Response" class="mb-2">
            <template #body>
                <record-item 
                    :titleGroup="[record.response.code, $constants.HTTP_STATUS_MESSAGE[record.response.code], record.response.proto]" 
                    :item="record.response">
                </record-item>
            </template>
        </box>

        <request-editor ref="reqEditor" :originalRequest="record.request"></request-editor>
        <code-viewer ref='curlViewer' :type="'shellscript'" :code="record.request.curl"></code-viewer>
    `,
    inject: ['$constants'],
    props: { record: Object },

    methods: {
        openCurlViewer() {
            if (!this.record.request.curl) {
                this.$bus.trigger('record.curl', this.record)
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
            <li @click='tab=0'>
                <button class="tab" :class="{ active: tab === 0 }">Header</button>
            </li>
            <li v-if="hasBody()" @click='tab=1; $refs.codeEditor.forceUpdate()'>
                <button class="tab" :class="{ active: tab === 1 }">Body</button>
            </li>

            <li v-if="isPreviewSupported()" @click="tab=2">
                <button class="tab" :class="{ active: tab === 2 }">Preview</button>
            </li>
        </ul>

        <div class="tab-content">
            <div class="tab-pane show" :class="{ active: tab === 0 }">
                <header-editor :header="item.header"></header-editor>
            </div>

            <div class="tab-pane mt-3" :class="{ active: tab === 1 }">
                <code-editor ref="codeEditor" :type="content.type" :code="item.body" :readOnly="true" :minLines="2" :maxLines="40"></code-editor>
            </div>

            <div class="tab-pane mt-3" :class="{ active: tab === 2 }">
                <base64-viewer :category="content.category" :type="content.type" :data="item.b64Body"></base64-viewer>
            </div>
        </div>
    `,
    inject: ['$util'],
    props: { titleGroup: Array, item: Object },
    data() { return { supportedPreviews: ['video', 'image'], tab: 0 } },
    beforeMount() { this.parseBody() },
    beforeUpdate() {
        this.parseBody()
        this.closeUnavailableTab()
    },
    methods: { 
        closeUnavailableTab() {
            let tab = this.tab
            if (tab === 2 && !this.isPreviewSupported()) { tab = 1 }
            if (tab === 1 && !this.hasBody()) { tab = 0 }
            this.tab = tab
        },

        parseBody() {
            this.content = {}
            if (this.item.body === void 0 || this.item.body === null) { this.item.body = '' }
            if (!this.item.header) { return }
            let contentType = this.item.header['Content-Type']
            if (!contentType) { return }
            this.content = this.$util.parseContentType(contentType.Val.join(";"))
        },

        hasBody() { return this.item.body !== '' },
        isPreviewSupported() { return this.supportedPreviews.indexOf(this.content.category) !== -1 }
    }
    
})