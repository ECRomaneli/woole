app.vue.component('Dashboard', {
    template:`
        <div class="container-fluid">
            <div class="row">
                <info-card class="col-xl-6 col-lg-8 col-md-12" :info="info"></info-card>
            </div>
        </div>
    `,
    props: { info: Object }
})

app.vue.component('InfoCard', {
    template: `
        <div class="p-0">
            <div class="card card-shadow" style="padding: 10px 20px">
                <div class="d-flex pt-3 pb-2"><span class="h2">Info</span></div>
                <table class="table table-striped table-hover info-grid" aria-label="info">
                    <tbody>
                        <tr v-for="(value, key) in info">
                            <th>{{ key }}</th>
                            <td v-html="getValue(value)"></td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    `,
    props: { info: Object },
    methods: {
        getValue(value) {
            if ((value + "").indexOf("://") !== -1) {
                return '<a target="_blank" href="' + value + '">' + value + '</a>'
            }

            return value
        }
    }
})

app.vue.component('PieChart', {
    template: `<canvas :id="id"></canvas>`,
    props: {
        rawData: Array,
        colorByLabel: Object

    },
    data() {
        return {
            id: 'pie-chart-' + app.nextInt(),
            defaultColors: [
                "#D32F2F", "#C2185B", "#7B1FA2", "#512DA8", "#303F9F", 
                "#1976D2", "#0288D1", "#0097A7", "#00796B", "#388E3C", 
                "#689F38", "#AFB42B", "#FBC02D", "#FFA000", "#F57C00", 
                "#E64A19", "#5D4037", "#616161", "#455A64"
            ]
        }
    },
    
    beforeUpdate() { updateData() },

    methods: {
        updateData() {
            const dataByLabel = {}

            this.label = []
            this.data = []

            this.rawData.forEach((rawData) => {
                if (!dataByLabel[rawData]) { dataByLabel[rawData] = 0 }
                dataByLabel[rawData]++
            })

            for (const label in data) {
                this.label.push(label)
                this.data.push(dataByLabel[label])
            }

            this.color = this.getColors(this.label)
        },

        getColors(labels) {            
            if (!this.colorByLabel) { return this.generateColors(labels.length) }
            return labels.map((l) => this.colorByLabel[l])
        },

        generateColors(number) {
            return this.defaultColors.slice(0, number)
        }
    }
})

// function getChartData(rawData) {
//     let dataByLabel = {}

//     let chart = { label: [], data: [], color: [] }

//     rawData.forEach((raw) => {
//         if (!dataByLabel[raw]) { dataByLabel[raw] = 0 }
//         dataByLabel[raw]++
//     })

//     for (const label in data) {
//         chart.label.push(label)
//         chart.data.push(dataByLabel[label])
//     }


// }