_N_E=(window.webpackJsonp_N_E=window.webpackJsonp_N_E||[]).push([[7],{"0sNQ":function(e,t){"trimStart"in String.prototype||(String.prototype.trimStart=String.prototype.trimLeft),"trimEnd"in String.prototype||(String.prototype.trimEnd=String.prototype.trimRight),"description"in Symbol.prototype||Object.defineProperty(Symbol.prototype,"description",{get:function(){return/\((.+)\)/.exec(this)[1]}}),Array.prototype.flat||(Array.prototype.flat=function(e,t){return t=this.concat.apply([],this),e>1&&t.some(Array.isArray)?t.flat(e-1):t},Array.prototype.flatMap=function(e,t){return this.map(e,t).flat()}),Promise.prototype.finally||(Promise.prototype.finally=function(e){if("function"!=typeof e)return this.then(e,e);var t=this.constructor||Promise;return this.then((function(r){return t.resolve(e()).then((function(){return r}))}),(function(r){return t.resolve(e()).then((function(){throw r}))}))})},"1ccW":function(e,t){function r(){return e.exports=r=Object.assign||function(e){for(var t=1;t<arguments.length;t++){var r=arguments[t];for(var n in r)Object.prototype.hasOwnProperty.call(r,n)&&(e[n]=r[n])}return e},r.apply(this,arguments)}e.exports=r},"48fX":function(e,t,r){var n=r("qhzo");e.exports=function(e,t){if("function"!==typeof t&&null!==t)throw new TypeError("Super expression must either be null or a function");e.prototype=Object.create(t&&t.prototype,{constructor:{value:e,writable:!0,configurable:!0}}),t&&n(e,t)}},BMP1:function(e,t,r){"use strict";var n=r("7KCV")(r("IKlv"));window.next=n,(0,n.default)().catch(console.error)},DqTX:function(e,t,r){"use strict";var n=r("zoAU");t.__esModule=!0,t.default=function(e){var t=document.getElementsByTagName("head")[0],r=new Set(t.children);i(r,e.map((function(e){var t=n(e,2),r=t[0],a=t[1];return(0,o.createElement)(r,a)})),!1);var a=null;return{mountedInstances:new Set,updateHead:function(e){var t=a=Promise.resolve().then((function(){t===a&&(a=null,i(r,e,!0))}))}}};var o=r("q1tI"),a={acceptCharset:"accept-charset",className:"class",htmlFor:"for",httpEquiv:"http-equiv"};function i(e,t,r){var n=document.getElementsByTagName("head")[0],o=new Set(e);t.forEach((function(t){if("title"!==t.type){for(var r=function(e){var t=e.type,r=e.props,n=document.createElement(t);for(var o in r)if(r.hasOwnProperty(o)&&"children"!==o&&"dangerouslySetInnerHTML"!==o&&void 0!==r[o]){var i=a[o]||o.toLowerCase();n.setAttribute(i,r[o])}var c=r.children,u=r.dangerouslySetInnerHTML;return u?n.innerHTML=u.__html||"":c&&(n.textContent="string"===typeof c?c:Array.isArray(c)?c.join(""):""),n}(t),i=e.values();;){var c=i.next(),u=c.done,s=c.value;if(null==s?void 0:s.isEqualNode(r))return void o.delete(s);if(u)break}e.add(r),n.appendChild(r)}else{var l="";if(t){var f=t.props.children;l="string"===typeof f?f:Array.isArray(f)?f.join(""):""}l!==document.title&&(document.title=l)}})),o.forEach((function(t){r&&t.parentNode.removeChild(t),e.delete(t)}))}},FYa8:function(e,t,r){"use strict";var n;t.__esModule=!0,t.HeadManagerContext=void 0;var o=((n=r("q1tI"))&&n.__esModule?n:{default:n}).default.createContext({});t.HeadManagerContext=o},IKlv:function(e,t,r){"use strict";var n=r("qVT1"),o=r("/GRZ"),a=r("i2R6"),i=r("48fX"),c=r("tCBg"),u=r("T0f4"),s=r("zoAU");function l(){l=function(){return e};var e={},t=Object.prototype,r=t.hasOwnProperty,n=Object.defineProperty||function(e,t,r){e[t]=r.value},o="function"==typeof Symbol?Symbol:{},a=o.iterator||"@@iterator",i=o.asyncIterator||"@@asyncIterator",c=o.toStringTag||"@@toStringTag";function u(e,t,r){return Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}),e[t]}try{u({},"")}catch(N){u=function(e,t,r){return e[t]=r}}function s(e,t,r,o){var a=t&&t.prototype instanceof p?t:p,i=Object.create(a.prototype),c=new T(o||[]);return n(i,"_invoke",{value:S(e,r,c)}),i}function f(e,t,r){try{return{type:"normal",arg:e.call(t,r)}}catch(N){return{type:"throw",arg:N}}}e.wrap=s;var d={};function p(){}function h(){}function v(){}var m={};u(m,a,(function(){return this}));var y=Object.getPrototypeOf,g=y&&y(y(L([])));g&&g!==t&&r.call(g,a)&&(m=g);var w=v.prototype=p.prototype=Object.create(m);function E(e){["next","throw","return"].forEach((function(t){u(e,t,(function(e){return this._invoke(t,e)}))}))}function _(e,t){var o;n(this,"_invoke",{value:function(n,a){function i(){return new t((function(o,i){!function n(o,a,i,c){var u=f(e[o],e,a);if("throw"!==u.type){var s=u.arg,l=s.value;return l&&"object"==typeof l&&r.call(l,"__await")?t.resolve(l.__await).then((function(e){n("next",e,i,c)}),(function(e){n("throw",e,i,c)})):t.resolve(l).then((function(e){s.value=e,i(s)}),(function(e){return n("throw",e,i,c)}))}c(u.arg)}(n,a,o,i)}))}return o=o?o.then(i,i):i()}})}function S(e,t,r){var n="suspendedStart";return function(o,a){if("executing"===n)throw new Error("Generator is already running");if("completed"===n){if("throw"===o)throw a;return C()}for(r.method=o,r.arg=a;;){var i=r.delegate;if(i){var c=b(i,r);if(c){if(c===d)continue;return c}}if("next"===r.method)r.sent=r._sent=r.arg;else if("throw"===r.method){if("suspendedStart"===n)throw n="completed",r.arg;r.dispatchException(r.arg)}else"return"===r.method&&r.abrupt("return",r.arg);n="executing";var u=f(e,t,r);if("normal"===u.type){if(n=r.done?"completed":"suspendedYield",u.arg===d)continue;return{value:u.arg,done:r.done}}"throw"===u.type&&(n="completed",r.method="throw",r.arg=u.arg)}}}function b(e,t){var r=t.method,n=e.iterator[r];if(void 0===n)return t.delegate=null,"throw"===r&&e.iterator.return&&(t.method="return",t.arg=void 0,b(e,t),"throw"===t.method)||"return"!==r&&(t.method="throw",t.arg=new TypeError("The iterator does not provide a '"+r+"' method")),d;var o=f(n,e.iterator,t.arg);if("throw"===o.type)return t.method="throw",t.arg=o.arg,t.delegate=null,d;var a=o.arg;return a?a.done?(t[e.resultName]=a.value,t.next=e.nextLoc,"return"!==t.method&&(t.method="next",t.arg=void 0),t.delegate=null,d):a:(t.method="throw",t.arg=new TypeError("iterator result is not an object"),t.delegate=null,d)}function x(e){var t={tryLoc:e[0]};1 in e&&(t.catchLoc=e[1]),2 in e&&(t.finallyLoc=e[2],t.afterLoc=e[3]),this.tryEntries.push(t)}function P(e){var t=e.completion||{};t.type="normal",delete t.arg,e.completion=t}function T(e){this.tryEntries=[{tryLoc:"root"}],e.forEach(x,this),this.reset(!0)}function L(e){if(e){var t=e[a];if(t)return t.call(e);if("function"==typeof e.next)return e;if(!isNaN(e.length)){var n=-1,o=function t(){for(;++n<e.length;)if(r.call(e,n))return t.value=e[n],t.done=!1,t;return t.value=void 0,t.done=!0,t};return o.next=o}}return{next:C}}function C(){return{value:void 0,done:!0}}return h.prototype=v,n(w,"constructor",{value:v,configurable:!0}),n(v,"constructor",{value:h,configurable:!0}),h.displayName=u(v,c,"GeneratorFunction"),e.isGeneratorFunction=function(e){var t="function"==typeof e&&e.constructor;return!!t&&(t===h||"GeneratorFunction"===(t.displayName||t.name))},e.mark=function(e){return Object.setPrototypeOf?Object.setPrototypeOf(e,v):(e.__proto__=v,u(e,c,"GeneratorFunction")),e.prototype=Object.create(w),e},e.awrap=function(e){return{__await:e}},E(_.prototype),u(_.prototype,i,(function(){return this})),e.AsyncIterator=_,e.async=function(t,r,n,o,a){void 0===a&&(a=Promise);var i=new _(s(t,r,n,o),a);return e.isGeneratorFunction(r)?i:i.next().then((function(e){return e.done?e.value:i.next()}))},E(w),u(w,c,"Generator"),u(w,a,(function(){return this})),u(w,"toString",(function(){return"[object Generator]"})),e.keys=function(e){var t=Object(e),r=[];for(var n in t)r.push(n);return r.reverse(),function e(){for(;r.length;){var n=r.pop();if(n in t)return e.value=n,e.done=!1,e}return e.done=!0,e}},e.values=L,T.prototype={constructor:T,reset:function(e){if(this.prev=0,this.next=0,this.sent=this._sent=void 0,this.done=!1,this.delegate=null,this.method="next",this.arg=void 0,this.tryEntries.forEach(P),!e)for(var t in this)"t"===t.charAt(0)&&r.call(this,t)&&!isNaN(+t.slice(1))&&(this[t]=void 0)},stop:function(){this.done=!0;var e=this.tryEntries[0].completion;if("throw"===e.type)throw e.arg;return this.rval},dispatchException:function(e){if(this.done)throw e;var t=this;function n(r,n){return i.type="throw",i.arg=e,t.next=r,n&&(t.method="next",t.arg=void 0),!!n}for(var o=this.tryEntries.length-1;o>=0;--o){var a=this.tryEntries[o],i=a.completion;if("root"===a.tryLoc)return n("end");if(a.tryLoc<=this.prev){var c=r.call(a,"catchLoc"),u=r.call(a,"finallyLoc");if(c&&u){if(this.prev<a.catchLoc)return n(a.catchLoc,!0);if(this.prev<a.finallyLoc)return n(a.finallyLoc)}else if(c){if(this.prev<a.catchLoc)return n(a.catchLoc,!0)}else{if(!u)throw new Error("try statement without catch or finally");if(this.prev<a.finallyLoc)return n(a.finallyLoc)}}}},abrupt:function(e,t){for(var n=this.tryEntries.length-1;n>=0;--n){var o=this.tryEntries[n];if(o.tryLoc<=this.prev&&r.call(o,"finallyLoc")&&this.prev<o.finallyLoc){var a=o;break}}a&&("break"===e||"continue"===e)&&a.tryLoc<=t&&t<=a.finallyLoc&&(a=null);var i=a?a.completion:{};return i.type=e,i.arg=t,a?(this.method="next",this.next=a.finallyLoc,d):this.complete(i)},complete:function(e,t){if("throw"===e.type)throw e.arg;return"break"===e.type||"continue"===e.type?this.next=e.arg:"return"===e.type?(this.rval=this.arg=e.arg,this.method="return",this.next="end"):"normal"===e.type&&t&&(this.next=t),d},finish:function(e){for(var t=this.tryEntries.length-1;t>=0;--t){var r=this.tryEntries[t];if(r.finallyLoc===e)return this.complete(r.completion,r.afterLoc),P(r),d}},catch:function(e){for(var t=this.tryEntries.length-1;t>=0;--t){var r=this.tryEntries[t];if(r.tryLoc===e){var n=r.completion;if("throw"===n.type){var o=n.arg;P(r)}return o}}throw new Error("illegal catch attempt")},delegateYield:function(e,t,r){return this.delegate={iterator:L(e),resultName:t,nextLoc:r},"next"===this.method&&(this.arg=void 0),d}},e}function f(e){var t=function(){if("undefined"===typeof Reflect||!Reflect.construct)return!1;if(Reflect.construct.sham)return!1;if("function"===typeof Proxy)return!0;try{return Boolean.prototype.valueOf.call(Reflect.construct(Boolean,[],(function(){}))),!0}catch(e){return!1}}();return function(){var r,n=u(e);if(t){var o=u(this).constructor;r=Reflect.construct(n,arguments,o)}else r=n.apply(this,arguments);return c(this,r)}}var d=r("7KCV"),p=r("AroE");t.__esModule=!0,t.render=ae,t.renderError=ce,t.default=t.emitter=t.router=t.version=void 0;var h=p(r("1ccW"));p(r("7KCV"));r("0sNQ");var v=p(r("q1tI")),m=p(r("i8i4")),y=r("FYa8"),g=p(r("dZ6Y")),w=r("qOIg"),E=r("elyg"),_=r("/jkW"),S=d(r("3WeD")),b=d(r("yLiY")),x=r("g/15"),P=p(r("DqTX")),T=d(r("zmvN")),L=p(r("bGXG")),C=r("nOHt"),N=JSON.parse(document.getElementById("__NEXT_DATA__").textContent);window.__NEXT_DATA__=N;t.version="9.5.5";var A=N.props,k=N.err,R=N.page,M=N.query,I=N.buildId,F=N.assetPrefix,j=N.runtimeConfig,O=N.dynamicIds,B=N.isFallback,D=N.head,q=N.locales,G=N.defaultLocale,H=N.locale,X=F||"";r.p="".concat(X,"/_next/"),b.setConfig({serverRuntimeConfig:{},publicRuntimeConfig:j||{}});var U=(0,x.getURL)();(0,E.hasBasePath)(U)&&(U=(0,E.delBasePath)(U)),U=(0,E.delLocale)(U,H);var W=new T.default(I,X,R),Y=function(e){var t=s(e,2),r=t[0],n=t[1];return W.registerPage(r,n)};window.__NEXT_P&&window.__NEXT_P.map((function(e){return setTimeout((function(){return Y(e)}),0)})),window.__NEXT_P=[],window.__NEXT_P.push=Y;var V,K,z,J,Z,Q,$,ee=(0,P.default)(D),te=document.getElementById("__next");t.router=z;var re=function(e){i(r,e);var t=f(r);function r(){return o(this,r),t.apply(this,arguments)}return a(r,[{key:"componentDidCatch",value:function(e,t){this.props.fn(e,t)}},{key:"componentDidMount",value:function(){this.scrollToHash(),z.isSsr&&(B||N.nextExport&&((0,_.isDynamicRoute)(z.pathname)||location.search)||A&&A.__N_SSG&&location.search)&&z.replace(z.pathname+"?"+String(S.assign(S.urlQueryToSearchParams(z.query),new URLSearchParams(location.search))),U,{_h:1,shallow:!B})}},{key:"componentDidUpdate",value:function(){this.scrollToHash()}},{key:"scrollToHash",value:function(){var e=location.hash;if(e=e&&e.substring(1)){var t=document.getElementById(e);t&&setTimeout((function(){return t.scrollIntoView()}),0)}}},{key:"render",value:function(){return this.props.children}}]),r}(v.default.Component),ne=(0,g.default)();t.emitter=ne;var oe=function(){var e=n(l().mark((function e(){var r,n,o,a,i,c,u=arguments;return l().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return u.length>0&&void 0!==u[0]?u[0]:{},e.next=4,W.loadPage("/_app");case 4:return r=e.sent,n=r.page,o=r.mod,Q=n,o&&o.reportWebVitals&&($=function(e){var t,r=e.id,n=e.name,a=e.startTime,i=e.value,c=e.duration,u=e.entryType,s=e.entries,l="".concat(Date.now(),"-").concat(Math.floor(8999999999999*Math.random())+1e12);s&&s.length&&(t=s[0].startTime),o.reportWebVitals({id:r||l,name:n,startTime:a||t,value:null==i?c:i,label:"mark"===u||"measure"===u?"custom":"web-vital"})}),a=k,e.prev=10,e.next=14,W.loadPage(R);case 14:i=e.sent,J=i.page,Z=i.styleSheets,e.next=21;break;case 21:e.next=26;break;case 23:e.prev=23,e.t0=e.catch(10),a=e.t0;case 26:if(!window.__NEXT_PRELOADREADY){e.next=30;break}return e.next=30,window.__NEXT_PRELOADREADY(O);case 30:return t.router=z=(0,C.createRouter)(R,M,U,{initialProps:A,pageLoader:W,App:Q,Component:J,initialStyleSheets:Z,wrapApp:pe,err:a,isFallback:Boolean(B),subscription:function(e,t){return ae({App:t,Component:e.Component,styleSheets:e.styleSheets,props:e.props,err:e.err})},locale:H,locales:q,defaultLocale:G}),ae(c={App:Q,Component:J,styleSheets:Z,props:A,err:a}),e.abrupt("return",ne);case 38:return e.abrupt("return",{emitter:ne,render:ae,renderCtx:c});case 39:case"end":return e.stop()}}),e,null,[[10,23]])})));return function(){return e.apply(this,arguments)}}();function ae(e){return ie.apply(this,arguments)}function ie(){return(ie=n(l().mark((function e(t){return l().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:if(!t.err){e.next=4;break}return e.next=3,ce(t);case 3:return e.abrupt("return");case 4:return e.prev=4,e.next=7,he(t);case 7:e.next=16;break;case 9:if(e.prev=9,e.t0=e.catch(4),!e.t0.cancelled){e.next=13;break}throw e.t0;case 13:return e.next=16,ce((0,h.default)({},t,{err:e.t0}));case 16:case"end":return e.stop()}}),e,null,[[4,9]])})))).apply(this,arguments)}function ce(e){var t=e.App,r=e.err;return console.error(r),W.loadPage("/_error").then((function(n){var o=n.page,a=n.styleSheets,i=pe(t),c={Component:o,AppTree:i,router:z,ctx:{err:r,pathname:R,query:M,asPath:U,AppTree:i}};return Promise.resolve(e.props?e.props:(0,x.loadGetInitialProps)(t,c)).then((function(t){return he((0,h.default)({},e,{err:r,Component:o,styleSheets:a,props:t}))}))}))}t.default=oe;var ue="function"===typeof m.default.hydrate;function se(){x.ST&&(performance.mark("afterHydrate"),performance.measure("Next.js-before-hydration","navigationStart","beforeRender"),performance.measure("Next.js-hydration","beforeRender","afterHydrate"),$&&performance.getEntriesByName("Next.js-hydration").forEach($),fe())}function le(){if(x.ST){performance.mark("afterRender");var e=performance.getEntriesByName("routeChange","mark");e.length&&(performance.measure("Next.js-route-change-to-render",e[0].name,"beforeRender"),performance.measure("Next.js-render","beforeRender","afterRender"),$&&(performance.getEntriesByName("Next.js-render").forEach($),performance.getEntriesByName("Next.js-route-change-to-render").forEach($)),fe(),["Next.js-route-change-to-render","Next.js-render"].forEach((function(e){return performance.clearMeasures(e)})))}}function fe(){["beforeRender","afterHydrate","afterRender","routeChange"].forEach((function(e){return performance.clearMarks(e)}))}function de(e){var t=e.children;return v.default.createElement(re,{fn:function(e){return ce({App:Q,err:e}).catch((function(e){return console.error("Error rendering page: ",e)}))}},v.default.createElement(w.RouterContext.Provider,{value:(0,C.makePublicRouterInstance)(z)},v.default.createElement(y.HeadManagerContext.Provider,{value:ee},t)))}var pe=function(e){return function(t){var r=(0,h.default)({},t,{Component:J,err:k,router:z});return v.default.createElement(de,null,v.default.createElement(e,r))}};function he(e){var t=e.App,r=e.Component,n=e.props,o=e.err,a=e.styleSheets;r=r||V.Component,n=n||V.props;var i=(0,h.default)({},n,{Component:r,err:o,router:z});V=i;var c,u=!1,s=new Promise((function(e,t){K&&K(),c=function(){K=null,e()},K=function(){u=!0,K=null;var e=new Error("Cancel rendering route");e.cancelled=!0,t(e)}}));var l,f,d=v.default.createElement(ve,{callback:function(){if(!ue&&!u){for(var e=new Set(a.map((function(e){return e.href}))),t=(0,T.looseToArray)(document.querySelectorAll("style[data-n-href]")),r=t.map((function(e){return e.getAttribute("data-n-href")})),n=0;n<r.length;++n)e.has(r[n])?t[n].removeAttribute("media"):t[n].setAttribute("media","x");var o=document.querySelector("noscript[data-n-css]");o&&a.forEach((function(e){var t=e.href,r=document.querySelector('style[data-n-href="'.concat(t,'"]'));r&&(o.parentNode.insertBefore(r,o.nextSibling),o=r)})),(0,T.looseToArray)(document.querySelectorAll("link[data-n-p]")).forEach((function(e){e.parentNode.removeChild(e)})),getComputedStyle(document.body,"height")}c()}},v.default.createElement(de,null,v.default.createElement(t,i)));return function(){if(ue)return!1;var e=(0,T.looseToArray)(document.querySelectorAll("style[data-n-href]")),t=new Set(e.map((function(e){return e.getAttribute("data-n-href")})));a.forEach((function(e){var r=e.href,n=e.text;if(!t.has(r)){var o=document.createElement("style");o.setAttribute("data-n-href",r),o.setAttribute("media","x"),document.head.appendChild(o),o.appendChild(document.createTextNode(n))}}))}(),l=d,f=te,x.ST&&performance.mark("beforeRender"),ue?(m.default.hydrate(l,f,se),ue=!1,$&&x.ST&&(0,L.default)($)):m.default.render(l,f,le),s}function ve(e){var t=e.callback,r=e.children;return v.default.useLayoutEffect((function(){return t()}),[t]),r}},Lab5:function(e,t,r){"use strict";t.__esModule=!0,t.default=function(e){var t=arguments.length>1&&void 0!==arguments[1]?arguments[1]:"",r="/"===e?"/index":/^\/index(\/|$)/.test(e)?"/index".concat(e):"".concat(e);return r+t}},T0f4:function(e,t){function r(t){return e.exports=r=Object.setPrototypeOf?Object.getPrototypeOf:function(e){return e.__proto__||Object.getPrototypeOf(e)},r(t)}e.exports=r},bGXG:function(e,t,r){"use strict";t.__esModule=!0,t.default=void 0;var n=r("w6Sm");t.default=function(e){(0,n.getCLS)(e),(0,n.getFID)(e),(0,n.getFCP)(e),(0,n.getLCP)(e),(0,n.getTTFB)(e)}},qXWd:function(e,t){e.exports=function(e){if(void 0===e)throw new ReferenceError("this hasn't been initialised - super() hasn't been called");return e}},tCBg:function(e,t,r){var n=r("C+bE"),o=r("qXWd");e.exports=function(e,t){return!t||"object"!==n(t)&&"function"!==typeof t?o(e):t}},w6Sm:function(e,t,r){"use strict";r.r(t),r.d(t,"getCLS",(function(){return h})),r.d(t,"getFCP",(function(){return m})),r.d(t,"getFID",(function(){return y})),r.d(t,"getLCP",(function(){return w})),r.d(t,"getTTFB",(function(){return E}));var n,o,a=function(){return"".concat(Date.now(),"-").concat(Math.floor(8999999999999*Math.random())+1e12)},i=function(e){var t=arguments.length>1&&void 0!==arguments[1]?arguments[1]:-1;return{name:e,value:t,delta:0,entries:[],id:a(),isFinal:!1}},c=function(e,t){try{if(PerformanceObserver.supportedEntryTypes.includes(e)){var r=new PerformanceObserver((function(e){return e.getEntries().map(t)}));return r.observe({type:e,buffered:!0}),r}}catch(e){}},u=!1,s=!1,l=function(e){u=!e.persisted},f=function(){addEventListener("pagehide",l),addEventListener("beforeunload",(function(){}))},d=function(e){var t=arguments.length>1&&void 0!==arguments[1]&&arguments[1];s||(f(),s=!0),addEventListener("visibilitychange",(function(t){var r=t.timeStamp;"hidden"===document.visibilityState&&e({timeStamp:r,isUnloading:u})}),{capture:!0,once:t})},p=function(e,t,r,n){var o;return function(){r&&t.isFinal&&r.disconnect(),t.value>=0&&(n||t.isFinal||"hidden"===document.visibilityState)&&(t.delta=t.value-(o||0),(t.delta||t.isFinal||void 0===o)&&(e(t),o=t.value))}},h=function(e){var t,r=arguments.length>1&&void 0!==arguments[1]&&arguments[1],n=i("CLS",0),o=function(e){e.hadRecentInput||(n.value+=e.value,n.entries.push(e),t())},a=c("layout-shift",o);a&&(t=p(e,n,a,r),d((function(e){var r=e.isUnloading;a.takeRecords().map(o),r&&(n.isFinal=!0),t()})))},v=function(){return void 0===n&&(n="hidden"===document.visibilityState?0:1/0,d((function(e){var t=e.timeStamp;return n=t}),!0)),{get timeStamp(){return n}}},m=function(e){var t,r=i("FCP"),n=v(),o=c("paint",(function(e){"first-contentful-paint"===e.name&&e.startTime<n.timeStamp&&(r.value=e.startTime,r.isFinal=!0,r.entries.push(e),t())}));o&&(t=p(e,r,o))},y=function(e){var t=i("FID"),r=v(),n=function(e){e.startTime<r.timeStamp&&(t.value=e.processingStart-e.startTime,t.entries.push(e),t.isFinal=!0,a())},o=c("first-input",n),a=p(e,t,o);o?d((function(){o.takeRecords().map(n),o.disconnect()}),!0):window.perfMetrics&&window.perfMetrics.onFirstInputDelay&&window.perfMetrics.onFirstInputDelay((function(e,n){n.timeStamp<r.timeStamp&&(t.value=e,t.isFinal=!0,t.entries=[{entryType:"first-input",name:n.type,target:n.target,cancelable:n.cancelable,startTime:n.timeStamp,processingStart:n.timeStamp+e}],a())}))},g=function(){return o||(o=new Promise((function(e){return["scroll","keydown","pointerdown"].map((function(t){addEventListener(t,e,{once:!0,passive:!0,capture:!0})}))}))),o},w=function(e){var t,r=arguments.length>1&&void 0!==arguments[1]&&arguments[1],n=i("LCP"),o=v(),a=function(e){var r=e.startTime;r<o.timeStamp?(n.value=r,n.entries.push(e)):n.isFinal=!0,t()},u=c("largest-contentful-paint",a);if(u){t=p(e,n,u,r);var s=function(){n.isFinal||(u.takeRecords().map(a),n.isFinal=!0,t())};g().then(s),d(s,!0)}},E=function(e){var t,r=i("TTFB");t=function(){try{var t=performance.getEntriesByType("navigation")[0]||function(){var e=performance.timing,t={entryType:"navigation",startTime:0};for(var r in e)"navigationStart"!==r&&"toJSON"!==r&&(t[r]=Math.max(e[r]-e.navigationStart,0));return t}();r.value=r.delta=t.responseStart,r.entries=[t],r.isFinal=!0,e(r)}catch(e){}},"complete"===document.readyState?setTimeout(t,0):addEventListener("pageshow",t)}},yLiY:function(e,t,r){"use strict";var n;t.__esModule=!0,t.setConfig=function(e){n=e},t.default=void 0;t.default=function(){return n}},zmvN:function(e,t,r){"use strict";var n=r("/GRZ"),o=r("i2R6"),a=r("AroE");t.__esModule=!0,t.default=t.looseToArray=void 0;var i=a(r("dZ6Y")),c=r("elyg"),u=a(r("Lab5")),s=r("/jkW"),l=r("hS4m"),f=function(e){return[].slice.call(e)};function d(e,t){try{return document.createElement("link").relList.supports(e)}catch(r){}}function p(e){return(0,c.markLoadingError)(new Error("Error loading ".concat(e)))}t.looseToArray=f;var h=d("preload")&&!d("prefetch")?"preload":"prefetch",v=d("preload")?"preload":h;document.createElement("script");function m(e){if("/"!==e[0])throw new Error('Route name should start with a "/", got "'.concat(e,'"'));return"/"===e?e:e.replace(/\/$/,"")}function y(e,t,r,n){return new Promise((function(o,a){n=document.createElement("link"),r&&(n.as=r),n.rel=t,n.crossOrigin=void 0,n.onload=o,n.onerror=a,n.href=e,document.head.appendChild(n)}))}var g=function(){function e(t,r,o){n(this,e),this.initialPage=void 0,this.buildId=void 0,this.assetPrefix=void 0,this.pageCache=void 0,this.pageRegisterEvents=void 0,this.loadingRoutes=void 0,this.promisedBuildManifest=void 0,this.promisedSsgManifest=void 0,this.promisedDevPagesManifest=void 0,this.initialPage=o,this.buildId=t,this.assetPrefix=r,this.pageCache={},this.pageRegisterEvents=(0,i.default)(),this.loadingRoutes={"/_app":!0},"/_error"!==o&&(this.loadingRoutes[o]=!0),this.promisedBuildManifest=new Promise((function(e){window.__BUILD_MANIFEST?e(window.__BUILD_MANIFEST):window.__BUILD_MANIFEST_CB=function(){e(window.__BUILD_MANIFEST)}})),this.promisedSsgManifest=new Promise((function(e){window.__SSG_MANIFEST?e(window.__SSG_MANIFEST):window.__SSG_MANIFEST_CB=function(){e(window.__SSG_MANIFEST)}}))}return o(e,[{key:"getPageList",value:function(){return this.promisedBuildManifest.then((function(e){return e.sortedPages}))}},{key:"getDependencies",value:function(e){var t=this;return this.promisedBuildManifest.then((function(r){return r[e]?r[e].map((function(e){return"".concat(t.assetPrefix,"/_next/").concat(encodeURI(e))})):Promise.reject(p(e))}))}},{key:"getDataHref",value:function(e,t,r,n,o){var a=this,i=(0,l.parseRelativeUrl)(e),f=i.pathname,d=i.query,p=i.search,h=(0,l.parseRelativeUrl)(t).pathname,v=m(f),y=function(e){var t=(0,c.addLocale)((0,u.default)(e,".json"),n,o);return(0,c.addBasePath)("/_next/data/".concat(a.buildId).concat(t).concat(r?"":p))},g=(0,s.isDynamicRoute)(v),w=g?(0,c.interpolateAs)(f,h,d).result:"";return g?w&&y(w):y(v)}},{key:"prefetchData",value:function(e,t,r,n){var o=this,a=m((0,l.parseRelativeUrl)(e).pathname);return this.promisedSsgManifest.then((function(i,c){return i.has(a)&&(c=o.getDataHref(e,t,!0,r,n))&&!document.querySelector('link[rel="'.concat(h,'"][href^="').concat(c,'"]'))&&y(c,h,"fetch").catch((function(){}))}))}},{key:"loadPage",value:function(e){var t=this;return e=m(e),new Promise((function(r,n){var o=t.pageCache[e];if(o)"error"in o?n(o.error):r(o);else{var a=function o(a){t.pageRegisterEvents.off(e,o),delete t.loadingRoutes[e],"error"in a?n(a.error):r(a)};if(t.pageRegisterEvents.on(e,a),!t.loadingRoutes[e])t.loadingRoutes[e]=!0,t.getDependencies(e).then((function(e){var t=[];return e.forEach((function(e){e.endsWith(".js")&&!document.querySelector('script[src^="'.concat(e,'"]'))&&t.push(function(e){return new Promise((function(t,r){var n=document.createElement("script");n.crossOrigin=void 0,n.src=e,n.onload=t,n.onerror=function(){return r(p(e))},document.body.appendChild(n)}))}(e)),e.endsWith(".css")&&!document.querySelector('link[rel="'.concat(v,'"][href^="').concat(e,'"]'))&&y(e,v,"fetch").catch((function(){}))})),Promise.all(t)})).catch((function(r){t.pageCache[e]={error:r},a({error:r})}))}}))}},{key:"registerPage",value:function(e,t){var r=this;var n=e===this.initialPage;("/_app"===e?Promise.resolve([]):(n?Promise.resolve(f(document.querySelectorAll("link[data-n-p]")).map((function(e){return e.getAttribute("href")}))):this.getDependencies(e).then((function(e){return e.filter((function(e){return e.endsWith(".css")}))}))).then((function(e){return Promise.all(e.map((function(e){return t=e,fetch(t).then((function(e){if(!e.ok)throw p(t);return e.text().then((function(e){return{href:t,text:e}}))}));var t}))).catch((function(e){if(n)return f(document.styleSheets).filter((function(e){return e.ownerNode&&"LINK"===e.ownerNode.tagName&&e.ownerNode.hasAttribute("data-n-p")})).map((function(e){return{href:e.ownerNode.getAttribute("href"),text:f(e.cssRules).map((function(e){return e.cssText})).join("")}}));throw e}))}))).then((function(n){return function(n){try{var o=t(),a={page:o.default||o,mod:o,styleSheets:n};r.pageCache[e]=a,r.pageRegisterEvents.emit(e,a)}catch(i){r.pageCache[e]={error:i},r.pageRegisterEvents.emit(e,{error:i})}}(n)}),(function(t){r.pageCache[e]={error:t},r.pageRegisterEvents.emit(e,{error:t})}))}},{key:"prefetch",value:function(e,t){var r,n,o=this;if((r=navigator.connection)&&(r.saveData||/2g/.test(r.effectiveType)))return Promise.resolve();if(t)n=e;else;return Promise.all(document.querySelector('link[rel="'.concat(h,'"][href^="').concat(n,'"]'))?[]:[n&&y(n,h,n.endsWith(".css")?"fetch":"script"),!t&&this.getDependencies(e).then((function(e){return Promise.all(e.map((function(e){return o.prefetch(e,!0)})))}))]).then((function(){}),(function(){}))}}]),e}();t.default=g}},[["BMP1",0,1,3]]]);