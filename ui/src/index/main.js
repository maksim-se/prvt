// This file is purposely written in ES5 style
/* eslint-disable */

;(function() {
    // Define all tests
    var features = {
        es6: function() {
            try {
                /*! ES6 detection credits DaBs and Netflix: https://gist.github.com/DaBs/89ccc2ffd1d435efdacff05248514f38 License: MIT */
                eval('class ಠ_ಠ extends Array {constructor(j = "a", ...c) {const q = (({u: e}) => {return { [`s${c}`]: Symbol(j) };})({});super(j, q, ...c);}}' + 
                'new Promise((f) => {const a = function* (){return "\\u{20BB7}".match(/./u)[0].length === 2 || true;};for (let vre of a()) {' +
                'const [uw, as, he, re] = [new Set(), new WeakSet(), new Map(), new WeakMap()];break;}f(new Proxy({}, {get: (han, h) => h in han ? han[h] ' + 
                ': "42".repeat(0o10)}));}).then(bi => new ಠ_ಠ(bi.rd));')
                return true
            }
            catch(e) {
                return false
            }
        },

        fetch: function() {
            return !!window.fetch
        },

        flexbox: function() {
            try {
                var element = document.createElement('b')
                element.style.display = 'flex'
                return element.style.display == 'flex'
            }
            catch (exception) {
                return false
            }
        },

        cssVariables: function() {
            return window.CSS && CSS.supports('color', 'var(--fake-var)')
        },

        cssTransitions: function() {
            return ('transition' in document.documentElement.style)
        }
    }

    // Run the tests
    var browserSupported = true
    for (var key in features) {
        if (features.hasOwnProperty(key)) {
            if (!features[key]()) {
                browserSupported = false
                break
            }
        }
    }

    if (browserSupported) {
        window.location.replace('app.html')
    } else {
        document.body.innerText = 'Your browser is not supported'
    }
})()
