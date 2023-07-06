app.component('Base64Viewer', {
    template: /*html*/ `
        <div class='base64-viewer'>
            <img v-if="category === 'image'" :src="contentToSource()" alt="preview" />
            <video :key="seqKey" v-else-if="category === 'video'" controls=""><source :type="getContentType()" :src="contentToSource()">Not Supported</video>
            <audio :key="seqKey" v-else-if="category === 'audio'" controls=""><source :type="getContentType()" :src="contentToSource()">Not supported</audio>
        </div>
    `,
    props: { category: String, type: String, data: String },
    data() { return { seqKey: 0 } },
    watch: {
        data() { this.seqKey++ }
    },
    methods: {
        supports(category) {
            return this.supportedCategories.some(c => c === category)
        },

        getContentType() {
            return this.category + '/' + this.type
        },
        // The body is already in base64
        contentToSource() {
            return "data:" + this.getContentType() + ";base64," + this.data
        }
    }
})