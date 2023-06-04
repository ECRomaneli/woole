app.component('RequestEditor', {
    template: /*html*/ `
        <modal ref="modal" @show="prepareRequest()" @hide="removeRequest()">
            <template #title>
                <img class="svg-icon square-24 me-2" :src="$image.src('request')" alt="request">
                <span class="h5 m-0">Request</span>
            </template>
            
            <template #body v-if="request">
                <div class="highlighted-group input-group">
                    <select class="request-method input-group-text" name="method" v-model="request.method">
                        <option v-for="(method) in $constants.HTTP_METHODS" :value="method">{{ method }}</option>
                    </select>
                    <input name="url" type="text" class="form-control" spellcheck="false" aria-label="url" v-model="request.url">
                </div>
                
                <div class="h5 centered-title mt-4 mb-3">Headers</div>
                <header-editor ref="headerEditor" :header="request.header" :read-only="false" @update="updateContent" @remove="removeContent"></header-editor>
                <div class="ps-2"><input type="checkbox" v-model="isAutoContentLengthEnabled"><label class="small-label">Calculate "Content-Length" on submit</label></div>
                
                <div class="h5 centered-title mt-4 mb-3">Body</div>
                <code-editor ref="codeEditor" :code="request.body" :type="content.type" :readOnly="false" :minLines="20" :maxLines="40"></code-editor>
            </template>

            <template #footer v-if="request">
                <button type="button" class="btn btn-sm" @click="$refs.modal.hide()">Cancel</button>
                <button type="button" class="btn btn-sm" @click="submit()">Submit</button>
            </template>
        </modal>
    `,
    inject: [ '$util', '$clone', '$constants', '$image' ],
    props: { originalRequest: Object },

    data() {
        return {
            content: {},
            request: null,
            isAutoContentLengthEnabled: true
        }
    },

    methods: {
        async submit() {
            this.request.b64Body = this.$refs.codeEditor.getCode()
            this.request.path = new URL(this.request.url, "http://dummy").pathname
            this.request.header = this.$refs.headerEditor.getHeader()

            if (this.isAutoContentLengthEnabled) {
                this.request.header['Content-Length'] = { Val: [`${this.$refs.codeEditor.getLength()}`] }
            }

            this.$bus.trigger('record.new', { request: this.request })
            this.$refs.modal.hide()
        },

        prepareRequest() {
            if (!this.originalRequest) {
                this.request = { method: "GET", header: [], proto: "HTTP/1.1"  }
                return
            }

            this.request = this.$clone(this.originalRequest)

            let contentType = this.request.header['Content-Type']
            this.updateContent({ name: "Content-Type", value: contentType ? contentType.Val.join(";") : "" })
        },

        removeRequest() { this.request = null },

        updateContent(header) {
            if (header.name.toLowerCase() === 'content-type') {
                this.content = this.$util.parseContentType(header.value)
            }
        },

        removeContent(header) {
            if (header.name.toLowerCase() === 'content-type') {
                this.content = {}
            }
        },

        show() {
            this.$refs.modal.show()
        }
    }
})