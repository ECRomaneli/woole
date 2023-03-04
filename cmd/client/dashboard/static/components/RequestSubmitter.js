app.vue.component('RequestSubmitter', {
    template: `
    <form @submit.prevent="submit" class="checkout-form">
        <div :id="modalId" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1" aria-labelledby="staticBackdropLabel" aria-hidden="true">
            <div class="modal-dialog modal-dialog-scrollable" style="max-width: 1000px">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title">Request</h5>
                    </div>
                    <div class="modal-body">
                        <div class="input-group mb-3">
                            <select class="request-method input-group-text" name="method" v-model="request.method">
                                <option v-for="(method) in httpMethods" :value="method">{{ method }}</option>
                            </select>
                            <input name="url" type="text" class="form-control" aria-label="url" v-model="request.url">
                        </div>
                        <div class="h5 centered-title">Headers</div>
                        <table class="table table-striped table-hover header-grid" aria-label="headers">
                            <thead>
                                <tr><th scope="remove"></th><th scope="name">Name</th><th scope="value">Value</th></tr>
                            </thead>
                            <tbody>
                                <tr v-for="(header, index) in request.header" :key="index">
                                    <td><div class="clickable-img" @click="remove(index)"><img class="bi" src="assets/images/trash.svg" :alt="'remove-header-'+index" style="width: 24px"></div></td>
                                    <td><textarea placeholder="Name" class="auto-resize" @focus="autoResize" @input="autoResize" @blur="autoResize" v-model="header.name"></textarea></td>
                                    <td><textarea placeholder="Value" class="auto-resize" @focus="autoResize" @input="autoResize" @blur="autoResize" v-model="header.value"></textarea></td>
                                </tr>
                                <tr>
                                    <td><div class="clickable-img" @click="add()"><img class="bi" src="assets/images/plus.svg" alt="add-header" style="width: 24px"></div></td>
                                    <td></td>
                                    <td></td>
                                </tr>
                            </tbody>
                        </table>
                        <div class="h5 centered-title">Body</div>
                        <content-editor ref="bodyEditor" :content="content" :readOnly="false" :minLines="20" :maxLines="40"></content-editor>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal" @click="cancel()">Cancel</button>
                        <button type="submit" class="btn btn-secondary">Submit</button>
                    </div>
                </div>
            </div>
        </div>
    </form>
    `,
    props: { modalId: String, originalRequest: Object },
    data() {
        return {
            httpMethods: ["HEAD", "GET", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"],
            content: { data: this.originalRequest.body }
        }
},
    beforeMount() { this.resetRequest() },
    watch: { originalRequest: function() { this.resetRequest() } },
    methods: {
        cancel() {
            this.resetRequest();
            this.$forceUpdate();
        },

        submit() {
            console.log("submitting...");
            this.request.body = this.$refs.bodyEditor.getValue();
            console.log(this.request);
        },

        resetRequest() {
            this.request = JSON.parse(JSON.stringify(this.originalRequest));
            this.request.header = [];

            Object.keys(this.originalRequest.header).forEach(headerName => {
                this.request.header.push({
                    name: headerName,
                    value: this.originalRequest.header[headerName].Val[0]
                });
            });
        },

        add() {
            this.request.header.push({ name: '', value: '' });
            this.$forceUpdate();
        },

        remove(index) {
            this.request.header.splice(index, 1);
            this.$forceUpdate();
        },

        autoResize(event) {
            let el = event.currentTarget;
            el.style.height = 'auto';
            el.style.height = event.type !== "blur" ? (el.scrollHeight)+"px" : "";
        }
    }
})