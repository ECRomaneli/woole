app.component('Base64Viewer', {
    template: /*html*/ `
        <div class='base64-viewer'>
            <img v-if="category === 'image'" :src="contentToSource()" alt="preview" />
            <video v-else-if="category === 'video'" controls=""><source :type="getContentType()" :src="contentToSource()"></video>
        </div>
    `,
    props: { category: String, type: String, data: String },
    methods: {
        getContentType() {
            return this.category + '/' + this.type
        },
        // The body is already in base64
        contentToSource() {
            return "data:" + this.getContentType() + ";base64," + this.data
        }
    }
})