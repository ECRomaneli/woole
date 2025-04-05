app.use({
    install: (app) => {
        const WEEKDAYS = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
        const MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

        class CustomDate {
            constructor (timestamp) {
                if (timestamp.toString().length === 10) { timestamp *= 1000; }
                this.date = new Date(timestamp)
            }

            format(format) {
                return format.replace(/(\w+)/g, (match) => {
                    switch (match) {
                        case 'ddd':  return WEEKDAYS[this.date.getDay()];
                        case 'DD':   return String(this.date.getDate()).padStart(2, '0');
                        case 'MMM':  return MONTHS[this.date.getMonth()];
                        case 'YYYY': return this.date.getFullYear();
                        case 'HH':   return String(this.date.getHours()).padStart(2, '0');
                        case 'hh':   return String(this.date.getHours() % 12 || 12).padStart(2, '0');
                        case 'A':    return this.date.getHours() < 12 ? 'AM' : 'PM';
                        case 'mm':   return String(this.date.getMinutes()).padStart(2, '0');
                        case 'ss':   return String(this.date.getSeconds()).padStart(2, '0');
                        case 'SSS':  return String(this.date.getMilliseconds()).padStart(3, '0');
                        default:     return match;
                    }
                });
            }
        }

        app.provide('$date', { from: (timestamp) => new CustomDate(timestamp) })
    }
})