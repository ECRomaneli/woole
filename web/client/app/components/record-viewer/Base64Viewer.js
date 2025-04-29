app.component('Base64Viewer', {
    template: /*html*/ `
        <div v-if='data' class='base64-viewer' @mouseover="showFitBtn = true" @mouseleave="showFitBtn = false">
            <img v-if="category === 'image'" :src="contentToSource()" :class="{ 'w-100 h-100': fitContainer }" alt="preview" />
            <video v-else-if="category === 'video' || category === 'application'" :key="seqKey" :class="{ 'w-100 h-100': fitContainer }" controls="">
                <source :type="getContentType()" :src="contentToSource()">Not Supported
            </video>
            <audio v-else-if="category === 'audio'" :key="seqKey" :class="{ 'w-100 h-100': fitContainer }" controls="">
                <source :type="getContentType()" :src="contentToSource()">Not Supported
            </audio>
            <Transition name="fast-fade">
                <button v-show="showFitBtn" class="btn" title="Toggle View" @click="fitContainer = !fitContainer">
                    <img class="svg-icon square-24" :src="$image.src(fitContainer ? 'minimize2' : 'maximize2')" alt="toggle-view">
                </button>
            </Transition>
        </div>
    `,
    inject: [ '$image' ],
    props: { category: String, type: String, data: String },
    data() { return { seqKey: 0, showFitBtn: false, fitContainer: false } },
    watch: {
        data() { this.seqKey++ }
    },
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