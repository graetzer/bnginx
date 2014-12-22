var Markdown;if(typeof exports==="object"&&typeof require==="function")Markdown=exports;else Markdown={};(function(){function e(e){return e}function t(e){return false}function n(){}function r(){}n.prototype={chain:function(t,n){var r=this[t];if(!r)throw new Error("unknown hook "+t);if(r===e)this[t]=n;else this[t]=function(e){var t=Array.prototype.slice.call(arguments,0);t[0]=r.apply(null,t);return n.apply(null,t)}},set:function(e,t){if(!this[e])throw new Error("unknown hook "+e);this[e]=t},addNoop:function(t){this[t]=e},addFalse:function(e){this[e]=t}};Markdown.HookCollection=n;r.prototype={set:function(e,t){this["s_"+e]=t},get:function(e){return this["s_"+e]}};Markdown.Converter=function(){function u(e){e=e.replace(/^[ ]{0,3}\[(.+)\]:[ \t]*\n?[ \t]*<?(\S+?)>?(?=\s|$)[ \t]*\n?[ \t]*((\n*)["(](.+?)[")][ \t]*)?(?:\n+)/gm,function(e,n,r,s,o,u){n=n.toLowerCase();t.set(n,O(r));if(o){return s}else if(u){i.set(n,u.replace(/"/g,"&quot;"))}return""});return e}function a(e){var t="p|div|h[1-6]|blockquote|pre|table|dl|ol|ul|script|noscript|form|fieldset|iframe|math|ins|del";var n="p|div|h[1-6]|blockquote|pre|table|dl|ol|ul|script|noscript|form|fieldset|iframe|math";e=e.replace(/^(<(p|div|h[1-6]|blockquote|pre|table|dl|ol|ul|script|noscript|form|fieldset|iframe|math|ins|del)\b[^\r]*?\n<\/\2>[ \t]*(?=\n+))/gm,f);e=e.replace(/^(<(p|div|h[1-6]|blockquote|pre|table|dl|ol|ul|script|noscript|form|fieldset|iframe|math)\b[^\r]*?.*<\/\2>[ \t]*(?=\n+)\n)/gm,f);e=e.replace(/\n[ ]{0,3}((<(hr)\b([^<>])*?\/?>)[ \t]*(?=\n{2,}))/g,f);e=e.replace(/\n\n[ ]{0,3}(<!(--(?:|(?:[^>-]|-[^>])(?:[^-]|-[^-])*)--)>[ \t]*(?=\n{2,}))/g,f);e=e.replace(/(?:\n\n)([ ]{0,3}(?:<([?%])[^\r]*?\2>)[ \t]*(?=\n{2,}))/g,f);return e}function f(e,t){var n=t;n=n.replace(/^\n+/,"");n=n.replace(/\n+$/g,"");n="\n\n~K"+(s.push(n)-1)+"K\n\n";return n}function c(t,n){t=e.preBlockGamut(t,l);t=b(t);var r="<hr />\n";t=t.replace(/^[ ]{0,2}([ ]?\*[ ]?){3,}[ \t]*$/gm,r);t=t.replace(/^[ ]{0,2}([ ]?-[ ]?){3,}[ \t]*$/gm,r);t=t.replace(/^[ ]{0,2}([ ]?_[ ]?){3,}[ \t]*$/gm,r);t=w(t);t=x(t);t=L(t);t=e.postBlockGamut(t,l);t=a(t);t=A(t,n);return t}function h(t){t=e.preSpanGamut(t);t=N(t);t=p(t);t=M(t);t=m(t);t=d(t);t=j(t);t=t.replace(/~P/g,"://");t=O(t);t=k(t);t=t.replace(/  +\n/g," <br>\n");t=e.postSpanGamut(t);return t}function p(e){var t=/(<[a-z\/!$]("[^"]*"|'[^']*'|[^'">])*>|<!(--(?:|(?:[^>-]|-[^>])(?:[^-]|-[^-])*)--)>)/gi;e=e.replace(t,function(e){var t=e.replace(/(.)<\/?code>(?=.)/g,"$1`");t=z(t,e.charAt(1)=="!"?"\\`*_/":"\\`*_");return t});return e}function d(e){e=e.replace(/(\[((?:\[[^\]]*\]|[^\[\]])*)\][ ]?(?:\n[ ]*)?\[(.*?)\])()()()()/g,v);e=e.replace(/(\[((?:\[[^\]]*\]|[^\[\]])*)\]\([ \t]*()<?((?:\([^)]*\)|[^()\s])*?)>?[ \t]*((['"])(.*?)\6[ \t]*)?\))/g,v);e=e.replace(/(\[([^\[\]]+)\])()()()()()/g,v);return e}function v(e,n,r,s,o,u,a,f){if(f==undefined)f="";var l=n;var c=r.replace(/:\/\//g,"~P");var h=s.toLowerCase();var p=o;var d=f;if(p==""){if(h==""){h=c.toLowerCase().replace(/ ?\n/g," ")}p="#"+h;if(t.get(h)!=undefined){p=t.get(h);if(i.get(h)!=undefined){d=i.get(h)}}else{if(l.search(/\(\s*\)$/m)>-1){p=""}else{return l}}}p=U(p);p=z(p,"*_");var v='<a href="'+p+'"';if(d!=""){d=g(d);d=z(d,"*_");v+=' title="'+d+'"'}v+=">"+c+"</a>";return v}function m(e){e=e.replace(/(!\[(.*?)\][ ]?(?:\n[ ]*)?\[(.*?)\])()()()()/g,y);e=e.replace(/(!\[(.*?)\]\s?\([ \t]*()<?(\S+?)>?[ \t]*((['"])(.*?)\6[ \t]*)?\))/g,y);return e}function g(e){return e.replace(/>/g,"&gt;").replace(/</g,"&lt;").replace(/"/g,"&quot;")}function y(e,n,r,s,o,u,a,f){var l=n;var c=r;var h=s.toLowerCase();var p=o;var d=f;if(!d)d="";if(p==""){if(h==""){h=c.toLowerCase().replace(/ ?\n/g," ")}p="#"+h;if(t.get(h)!=undefined){p=t.get(h);if(i.get(h)!=undefined){d=i.get(h)}}else{return l}}c=z(g(c),"*_[]()");p=z(p,"*_");var v='<img src="'+p+'" alt="'+c+'"';d=g(d);d=z(d,"*_");v+=' title="'+d+'"';v+=" />";return v}function b(e){e=e.replace(/^(.+)[ \t]*\n=+[ \t]*\n+/gm,function(e,t){return"<h1>"+h(t)+"</h1>\n\n"});e=e.replace(/^(.+)[ \t]*\n-+[ \t]*\n+/gm,function(e,t){return"<h2>"+h(t)+"</h2>\n\n"});e=e.replace(/^(\#{1,6})[ \t]*(.+?)[ \t]*\#*\n+/gm,function(e,t,n){var r=t.length;return"<h"+r+">"+h(n)+"</h"+r+">\n\n"});return e}function w(e,t){e+="~0";var n=/^(([ ]{0,3}([*+-]|\d+[.])[ \t]+)[^\r]+?(~0|\n{2,}(?=\S)(?![ \t]*(?:[*+-]|\d+[.])[ \t]+)))/gm;if(o){e=e.replace(n,function(e,n,r){var i=n;var s=r.search(/[*+-]/g)>-1?"ul":"ol";var o=S(i,s,t);o=o.replace(/\s+$/,"");o="<"+s+">"+o+"</"+s+">\n";return o})}else{n=/(\n\n|^\n?)(([ ]{0,3}([*+-]|\d+[.])[ \t]+)[^\r]+?(~0|\n{2,}(?=\S)(?![ \t]*(?:[*+-]|\d+[.])[ \t]+)))/g;e=e.replace(n,function(e,t,n,r){var i=t;var s=n;var o=r.search(/[*+-]/g)>-1?"ul":"ol";var u=S(s,o);u=i+"<"+o+">\n"+u+"</"+o+">\n";return u})}e=e.replace(/~0/,"");return e}function S(e,t,n){o++;e=e.replace(/\n{2,}$/,"\n");e+="~0";var r=E[t];var i=new RegExp("(^[ \\t]*)("+r+")[ \\t]+([^\\r]+?(\\n+))(?=(~0|\\1("+r+")[ \\t]+))","gm");var s=false;e=e.replace(i,function(e,t,r,i){var o=i;var u=t;var a=/\n\n$/.test(o);var f=a||o.search(/\n{2,}/)>-1;if(f||s){o=c(I(o),true)}else{o=w(I(o),true);o=o.replace(/\n$/,"");if(!n)o=h(o)}s=a;return"<li>"+o+"</li>\n"});e=e.replace(/~0/g,"");o--;return e}function x(e){e+="~0";e=e.replace(/(?:\n\n|^\n?)((?:(?:[ ]{4}|\t).*\n+)+)(\n*[ ]{0,3}[^ \t\n]|(?=~0))/g,function(e,t,n){var r=t;var i=n;r=C(I(r));r=q(r);r=r.replace(/^\n+/g,"");r=r.replace(/\n+$/g,"");r="<pre><code>"+r+"\n</code></pre>";return"\n\n"+r+"\n\n"+i});e=e.replace(/~0/,"");return e}function T(e){e=e.replace(/(^\n+|\n+$)/g,"");return"\n\n~K"+(s.push(e)-1)+"K\n\n"}function N(e){e=e.replace(/(^|[^\\])(`+)([^\r]*?[^`])\2(?!`)/gm,function(e,t,n,r,i){var s=r;s=s.replace(/^([ \t]*)/g,"");s=s.replace(/[ \t]*$/g,"");s=C(s);s=s.replace(/:\/\//g,"~P");return t+"<code>"+s+"</code>"});return e}function C(e){e=e.replace(/&/g,"&");e=e.replace(/</g,"&lt;");e=e.replace(/>/g,"&gt;");e=z(e,"*_{}[]\\",false);return e}function k(e){e=e.replace(/([\W_]|^)(\*\*|__)(?=\S)([^\r]*?\S[\*_]*)\2([\W_]|$)/g,"$1<strong>$3</strong>$4");e=e.replace(/([\W_]|^)(\*|_)(?=\S)([^\r\*_]*?\S)\2([\W_]|$)/g,"$1<em>$3</em>$4");return e}function L(e){e=e.replace(/((^[ \t]*>[ \t]?.+\n(.+\n)*\n*)+)/gm,function(e,t){var n=t;n=n.replace(/^[ \t]*>[ \t]?/gm,"~0");n=n.replace(/~0/g,"");n=n.replace(/^[ \t]+$/gm,"");n=c(n);n=n.replace(/(^|\n)/g,"$1  ");n=n.replace(/(\s*<pre>[^\r]+?<\/pre>)/gm,function(e,t){var n=t;n=n.replace(/^  /mg,"~0");n=n.replace(/~0/g,"");return n});return T("<blockquote>\n"+n+"\n</blockquote>")});return e}function A(e,t){e=e.replace(/^\n+/g,"");e=e.replace(/\n+$/g,"");var n=e.split(/\n{2,}/g);var r=[];var i=/~K(\d+)K/;var o=n.length;for(var u=0;u<o;u++){var a=n[u];if(i.test(a)){r.push(a)}else if(/\S/.test(a)){a=h(a);a=a.replace(/^([ \t]*)/g,"<p>");a+="</p>";r.push(a)}}if(!t){o=r.length;for(var u=0;u<o;u++){var f=true;while(f){f=false;r[u]=r[u].replace(/~K(\d+)K/g,function(e,t){f=true;return s[t]})}}}return r.join("\n\n")}function O(e){e=e.replace(/&(?!#?[xX]?(?:[0-9a-fA-F]+|\w+);)/g,"&");e=e.replace(/<(?![a-z\/?!]|~D)/gi,"&lt;");return e}function M(e){e=e.replace(/\\(\\)/g,W);e=e.replace(/\\([`*_{}\[\]()>#+-.!])/g,W);return e}function B(e,t,n,r){if(t)return e;if(r.charAt(r.length-1)!==")")return"<"+n+r+">";var i=r.match(/[()]/g);var s=0;for(var o=0;o<i.length;o++){if(i[o]==="("){if(s<=0)s=1;else s++}else{s--}}var u="";if(s<0){var a=new RegExp("\\){1,"+ -s+"}$");r=r.replace(a,function(e){u=e;return""})}if(u){var f=r.charAt(r.length-1);if(!H.test(f)){u=f+u;r=r.substr(0,r.length-1)}}return"<"+n+r+">"+u}function j(t){t=t.replace(P,B);var n=function(t,n){return'<a href="'+n+'">'+e.plainLinkText(n)+"</a>"};t=t.replace(/<((https?|ftp):[^'">\s]+)>/gi,n);return t}function F(e){e=e.replace(/~E(\d+)E/g,function(e,t){var n=parseInt(t);return String.fromCharCode(n)});return e}function I(e){e=e.replace(/^(\t|[ ]{1,4})/gm,"~0");e=e.replace(/~0/g,"");return e}function q(e){if(!/\t/.test(e))return e;var t=["    ","   ","  "," "],n=0,r;return e.replace(/[\n\t]/g,function(e,i){if(e==="\n"){n=i+1;return e}r=(i-n)%4;n=i+1;return t[r]})}function U(e){if(!e)return"";var t=e.length;return e.replace(R,function(n,r){if(n=="~D")return"%24";if(n==":"){if(r==t-1||/[0-9\/]/.test(e.charAt(r+1)))return":"}return"%"+n.charCodeAt(0).toString(16)})}function z(e,t,n){var r="(["+t.replace(/([\[\]\\])/g,"\\$1")+"])";if(n){r="\\\\"+r}var i=new RegExp(r,"g");e=e.replace(i,W);return e}function W(e,t){var n=t.charCodeAt(0);return"~E"+n+"E"}var e=this.hooks=new n;e.addNoop("plainLinkText");e.addNoop("preConversion");e.addNoop("postNormalization");e.addNoop("preBlockGamut");e.addNoop("postBlockGamut");e.addNoop("preSpanGamut");e.addNoop("postSpanGamut");e.addNoop("postConversion");var t;var i;var s;var o;this.makeHtml=function(n){if(t)throw new Error("Recursive call to converter.makeHtml");t=new r;i=new r;s=[];o=0;n=e.preConversion(n);n=n.replace(/~/g,"~T");n=n.replace(/\$/g,"~D");n=n.replace(/\r\n/g,"\n");n=n.replace(/\r/g,"\n");n="\n\n"+n+"\n\n";n=q(n);n=n.replace(/^[ \t]+$/mg,"");n=e.postNormalization(n);n=a(n);n=u(n);n=c(n);n=F(n);n=n.replace(/~D/g,"$$");n=n.replace(/~T/g,"~");n=e.postConversion(n);s=i=t=null;return n};var l=function(e){return c(e)};var E={ol:"\\d+[.]",ul:"[*+-]"};var _="[-A-Z0-9+&@#/%?=~_|[\\]()!:,.;]",D="[-A-Z0-9+&@#/%=~_|[\\])]",P=new RegExp('(="|<)?\\b(https?|ftp)(://'+_+"*"+D+")(?=$|\\W)","gi"),H=new RegExp(D,"i");var R=/(?:["'*()[\]:]|~D)/g}})()