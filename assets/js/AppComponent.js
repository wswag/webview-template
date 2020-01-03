var app = new Vue({
    el: '#app',
    data: {
        vm: model
    },
    methods: {
        hi: function() { model.hi() }
    }
});