_N_E=(window.webpackJsonp_N_E=window.webpackJsonp_N_E||[]).push([[14],{o4Gj:function(t,e,r){"use strict";r.r(e);var n=r("HaE+"),o=r("ODXe"),a=r("vDqi"),i=r.n(a),c=r("nOHt"),u=r.n(c),s=r("q1tI"),l=r.n(s),f=r("VtrM"),h=r("PSzF"),p=r("FBIo"),d=r("oR5U"),v=r("HPzk"),y=r("OZLi"),g=r("1EOh"),m=l.a.createElement;function w(){w=function(){return t};var t={},e=Object.prototype,r=e.hasOwnProperty,n=Object.defineProperty||function(t,e,r){t[e]=r.value},o="function"==typeof Symbol?Symbol:{},a=o.iterator||"@@iterator",i=o.asyncIterator||"@@asyncIterator",c=o.toStringTag||"@@toStringTag";function u(t,e,r){return Object.defineProperty(t,e,{value:r,enumerable:!0,configurable:!0,writable:!0}),t[e]}try{u({},"")}catch(k){u=function(t,e,r){return t[e]=r}}function s(t,e,r,o){var a=e&&e.prototype instanceof h?e:h,i=Object.create(a.prototype),c=new N(o||[]);return n(i,"_invoke",{value:E(t,r,c)}),i}function l(t,e,r){try{return{type:"normal",arg:t.call(e,r)}}catch(k){return{type:"throw",arg:k}}}t.wrap=s;var f={};function h(){}function p(){}function d(){}var v={};u(v,a,(function(){return this}));var y=Object.getPrototypeOf,g=y&&y(y(j([])));g&&g!==e&&r.call(g,a)&&(v=g);var m=d.prototype=h.prototype=Object.create(v);function b(t){["next","throw","return"].forEach((function(e){u(t,e,(function(t){return this._invoke(e,t)}))}))}function x(t,e){var o;n(this,"_invoke",{value:function(n,a){function i(){return new e((function(o,i){!function n(o,a,i,c){var u=l(t[o],t,a);if("throw"!==u.type){var s=u.arg,f=s.value;return f&&"object"==typeof f&&r.call(f,"__await")?e.resolve(f.__await).then((function(t){n("next",t,i,c)}),(function(t){n("throw",t,i,c)})):e.resolve(f).then((function(t){s.value=t,i(s)}),(function(t){return n("throw",t,i,c)}))}c(u.arg)}(n,a,o,i)}))}return o=o?o.then(i,i):i()}})}function E(t,e,r){var n="suspendedStart";return function(o,a){if("executing"===n)throw new Error("Generator is already running");if("completed"===n){if("throw"===o)throw a;return T()}for(r.method=o,r.arg=a;;){var i=r.delegate;if(i){var c=L(i,r);if(c){if(c===f)continue;return c}}if("next"===r.method)r.sent=r._sent=r.arg;else if("throw"===r.method){if("suspendedStart"===n)throw n="completed",r.arg;r.dispatchException(r.arg)}else"return"===r.method&&r.abrupt("return",r.arg);n="executing";var u=l(t,e,r);if("normal"===u.type){if(n=r.done?"completed":"suspendedYield",u.arg===f)continue;return{value:u.arg,done:r.done}}"throw"===u.type&&(n="completed",r.method="throw",r.arg=u.arg)}}}function L(t,e){var r=e.method,n=t.iterator[r];if(void 0===n)return e.delegate=null,"throw"===r&&t.iterator.return&&(e.method="return",e.arg=void 0,L(t,e),"throw"===e.method)||"return"!==r&&(e.method="throw",e.arg=new TypeError("The iterator does not provide a '"+r+"' method")),f;var o=l(n,t.iterator,e.arg);if("throw"===o.type)return e.method="throw",e.arg=o.arg,e.delegate=null,f;var a=o.arg;return a?a.done?(e[t.resultName]=a.value,e.next=t.nextLoc,"return"!==e.method&&(e.method="next",e.arg=void 0),e.delegate=null,f):a:(e.method="throw",e.arg=new TypeError("iterator result is not an object"),e.delegate=null,f)}function _(t){var e={tryLoc:t[0]};1 in t&&(e.catchLoc=t[1]),2 in t&&(e.finallyLoc=t[2],e.afterLoc=t[3]),this.tryEntries.push(e)}function O(t){var e=t.completion||{};e.type="normal",delete e.arg,t.completion=e}function N(t){this.tryEntries=[{tryLoc:"root"}],t.forEach(_,this),this.reset(!0)}function j(t){if(t){var e=t[a];if(e)return e.call(t);if("function"==typeof t.next)return t;if(!isNaN(t.length)){var n=-1,o=function e(){for(;++n<t.length;)if(r.call(t,n))return e.value=t[n],e.done=!1,e;return e.value=void 0,e.done=!0,e};return o.next=o}}return{next:T}}function T(){return{value:void 0,done:!0}}return p.prototype=d,n(m,"constructor",{value:d,configurable:!0}),n(d,"constructor",{value:p,configurable:!0}),p.displayName=u(d,c,"GeneratorFunction"),t.isGeneratorFunction=function(t){var e="function"==typeof t&&t.constructor;return!!e&&(e===p||"GeneratorFunction"===(e.displayName||e.name))},t.mark=function(t){return Object.setPrototypeOf?Object.setPrototypeOf(t,d):(t.__proto__=d,u(t,c,"GeneratorFunction")),t.prototype=Object.create(m),t},t.awrap=function(t){return{__await:t}},b(x.prototype),u(x.prototype,i,(function(){return this})),t.AsyncIterator=x,t.async=function(e,r,n,o,a){void 0===a&&(a=Promise);var i=new x(s(e,r,n,o),a);return t.isGeneratorFunction(r)?i:i.next().then((function(t){return t.done?t.value:i.next()}))},b(m),u(m,c,"Generator"),u(m,a,(function(){return this})),u(m,"toString",(function(){return"[object Generator]"})),t.keys=function(t){var e=Object(t),r=[];for(var n in e)r.push(n);return r.reverse(),function t(){for(;r.length;){var n=r.pop();if(n in e)return t.value=n,t.done=!1,t}return t.done=!0,t}},t.values=j,N.prototype={constructor:N,reset:function(t){if(this.prev=0,this.next=0,this.sent=this._sent=void 0,this.done=!1,this.delegate=null,this.method="next",this.arg=void 0,this.tryEntries.forEach(O),!t)for(var e in this)"t"===e.charAt(0)&&r.call(this,e)&&!isNaN(+e.slice(1))&&(this[e]=void 0)},stop:function(){this.done=!0;var t=this.tryEntries[0].completion;if("throw"===t.type)throw t.arg;return this.rval},dispatchException:function(t){if(this.done)throw t;var e=this;function n(r,n){return i.type="throw",i.arg=t,e.next=r,n&&(e.method="next",e.arg=void 0),!!n}for(var o=this.tryEntries.length-1;o>=0;--o){var a=this.tryEntries[o],i=a.completion;if("root"===a.tryLoc)return n("end");if(a.tryLoc<=this.prev){var c=r.call(a,"catchLoc"),u=r.call(a,"finallyLoc");if(c&&u){if(this.prev<a.catchLoc)return n(a.catchLoc,!0);if(this.prev<a.finallyLoc)return n(a.finallyLoc)}else if(c){if(this.prev<a.catchLoc)return n(a.catchLoc,!0)}else{if(!u)throw new Error("try statement without catch or finally");if(this.prev<a.finallyLoc)return n(a.finallyLoc)}}}},abrupt:function(t,e){for(var n=this.tryEntries.length-1;n>=0;--n){var o=this.tryEntries[n];if(o.tryLoc<=this.prev&&r.call(o,"finallyLoc")&&this.prev<o.finallyLoc){var a=o;break}}a&&("break"===t||"continue"===t)&&a.tryLoc<=e&&e<=a.finallyLoc&&(a=null);var i=a?a.completion:{};return i.type=t,i.arg=e,a?(this.method="next",this.next=a.finallyLoc,f):this.complete(i)},complete:function(t,e){if("throw"===t.type)throw t.arg;return"break"===t.type||"continue"===t.type?this.next=t.arg:"return"===t.type?(this.rval=this.arg=t.arg,this.method="return",this.next="end"):"normal"===t.type&&e&&(this.next=e),f},finish:function(t){for(var e=this.tryEntries.length-1;e>=0;--e){var r=this.tryEntries[e];if(r.finallyLoc===t)return this.complete(r.completion,r.afterLoc),O(r),f}},catch:function(t){for(var e=this.tryEntries.length-1;e>=0;--e){var r=this.tryEntries[e];if(r.tryLoc===t){var n=r.completion;if("throw"===n.type){var o=n.arg;O(r)}return o}}throw new Error("illegal catch attempt")},delegateYield:function(t,e,r){return this.delegate={iterator:j(t),resultName:e,nextLoc:r},"next"===this.method&&(this.arg=void 0),f}},t}var b=function(t){var e=t.article,r={title:e.title,description:e.description,body:e.body,tagList:e.tagList},a=l.a.useState(!1),s=Object(o.a)(a,2),d=s[0],b=s[1],x=l.a.useState([]),E=Object(o.a)(x,2),L=E[0],_=E[1],O=l.a.useReducer(y.a,r),N=Object(o.a)(O,2),j=N[0],T=N[1],k=Object(f.a)("user",g.a).data,S=Object(c.useRouter)().query.pid,P=function(){var t=Object(n.a)(w().mark((function t(e){var r,n,o;return w().wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return e.preventDefault(),b(!0),t.next=4,i.a.put("".concat(v.e,"/articles/").concat(S),JSON.stringify({article:j}),{headers:{"Content-Type":"application/json",Authorization:"Token ".concat(encodeURIComponent(null===k||void 0===k?void 0:k.token))}});case 4:r=t.sent,n=r.data,o=r.status,b(!1),200!==o&&_(n.errors),u.a.push("/");case 10:case"end":return t.stop()}}),t)})));return function(e){return t.apply(this,arguments)}}();return m("div",{className:"editor-page"},m("div",{className:"container page"},m("div",{className:"row"},m("div",{className:"col-md-10 offset-md-1 col-xs-12"},m(h.a,{errors:L}),m("form",null,m("fieldset",null,m("fieldset",{className:"form-group"},m("input",{className:"form-control form-control-lg",type:"text",placeholder:"Article Title",value:j.title,onChange:function(t){return T({type:"SET_TITLE",text:t.target.value})}})),m("fieldset",{className:"form-group"},m("input",{className:"form-control",type:"text",placeholder:"What's this article about?",value:j.description,onChange:function(t){return T({type:"SET_DESCRIPTION",text:t.target.value})}})),m("fieldset",{className:"form-group"},m("textarea",{className:"form-control",rows:8,placeholder:"Write your article (in markdown)",value:j.body,onChange:function(t){return T({type:"SET_BODY",text:t.target.value})}})),m(p.a,{tagList:j.tagList,addTag:function(t){return T({type:"ADD_TAG",tag:t})},removeTag:function(t){return T({type:"REMOVE_TAG",tag:t})}}),m("button",{className:"btn btn-lg pull-xs-right btn-primary",type:"button",disabled:d,onClick:P},"Update Article")))))))};b.getInitialProps=function(){var t=Object(n.a)(w().mark((function t(e){var r,n,o;return w().wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return r=e.query.pid,t.next=3,d.a.get(r);case 3:return n=t.sent,o=n.data.article,t.abrupt("return",{article:o});case 6:case"end":return t.stop()}}),t)})));return function(e){return t.apply(this,arguments)}}(),e.default=b},sU0j:function(t,e,r){(window.__NEXT_P=window.__NEXT_P||[]).push(["/editor/[pid]",function(){return r("o4Gj")}])}},[["sU0j",0,2,1,3,4,8]]]);