webpackJsonp([1],{24:function(e,t,s){"use strict";var n=s(2),a=s(87),i=s(82),o=s.n(i);n.default.use(a.a),t.a=new a.a({routes:[{path:"/",name:"Intro",component:o.a}]})},26:function(e,t){},27:function(e,t,s){s(79);var n=s(14)(s(52),s(86),null,null);e.exports=n.exports},51:function(e,t,s){"use strict";Object.defineProperty(t,"__esModule",{value:!0});var n=s(2),a=s(25),i=s.n(a),o=s(26),r=(s.n(o),s(27)),l=s.n(r),c=s(24),u=s(28);n.default.use(i.a),n.default.use(u.a),n.default.config.productionTip=!1,new n.default({el:"#app",router:c.a,template:"<App/>",components:{App:l.a}})},52:function(e,t,s){"use strict";Object.defineProperty(t,"__esModule",{value:!0}),t.default={name:"app"}},53:function(e,t,s){"use strict";Object.defineProperty(t,"__esModule",{value:!0});var n=s(83),a=s.n(n);t.default={name:"intro",components:{Try:a.a},data:function(){return{}}}},54:function(e,t,s){"use strict";Object.defineProperty(t,"__esModule",{value:!0});var n=s(55),a=s.n(n);t.default={name:"try",data:function(){return{testQueue:"test",testSquad:"gosqs.org",newMessage:"",messageLogs:[],masterAddr:"",servingNode:"",token:"",nextId:0}},created:function(){this.masterAddr="http://127.0.0.1:5446",this.applyNode()},mounted:function(){console.log("mounted"),setInterval(this.pullMessage,1e3)},computed:{apiJSONparams:function(){return a()({queue_name:this.testQueue,squad_name:this.testSquad,token:this.token,message_id:this.nextId,content:this.newMessage,size:1})},testingTagType:function(){return""===this.servingNode?"gray":"success"},testingAvailable:function(){return""===this.servingNode},testingStateString:function(){return""===this.servingNode?"unavailabe":"availabe"}},methods:{applyNode:function(){var e=this,t=a()({access_key:"Focinfi",secret_key:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyIiOiJGb2NpbmZpIn0.JRfcSY6_syXNKypblVpdD7oNmHNFDKVPVjkbNTlzEKQ",queue_name:this.testQueue,squad_name:this.testSquad});this.$http.post(this.masterAddr+"/applyNode",t).then(function(t){var s=t.body;if(1e3===s.code){var n=s.data;e.servingNode=n.node,e.token=n.token}},function(e){})},pullMessage:function(){var e=this;if(""===this.servingNode)return void(this.messageLogs=[this.timeFormat(new Date)+" Failed to connect the node"]);this.$http.post(this.servingNode+"/messages",this.apiJSONparams).then(function(t){var s=t.body.data,n=e.timeFormat(new Date),a={time:n,messages:["No More Messages"]};s.messages.length>0&&(a.messages=s.messages,e.reportReceived()),e.messageLogs.push(a),e.messageLogs.length>5&&e.messageLogs.shift()},function(e){console.log(e)})},pushMessage:function(){var e=this,t=this;this.$prompt("Input Message","Ready To Push",{confirmButtonText:"Push it",cancelButtonText:"Cancel",inputValidator:function(e){return null===e||void 0===e||""===e?"Can not push empty message":e.length>20?"Content is too long":(console.log(e),t.newMessage=e,!0)}}).then(function(s){s.value;e.applyMessageID(function(){t.$http.post(t.servingNode+"/message",t.apiJSONparams).then(function(e){t.newMessage="",console.log(e.data),t.$message({type:"success",message:"Message pushed"})},function(e){t.$message({type:"fail",message:"Failed to push message"})})})}).catch(function(){e.$message({type:"info",message:"Cancel Push"})})},applyMessageID:function(e){var t=this;this.$http.post(this.servingNode+"/messageID",this.apiJSONparams).then(function(s){console.log(s.data);var n=s.body.data;t.nextId=n.message_id_end,console.log("nextID",t.nextId),e()},function(e){})},reportReceived:function(){this.$http.post(this.servingNode+"/receivedMessageID",this.apiJSONparams).then(function(e){},function(e){console.log(e)})},messagesLogsSlice:function(e){console.log(this.messageLogs.length)},timeFormat:function(e){return"["+e.toLocaleTimeString("en-GB")+"]"},messagesEntryFormat:function(e){return"["+e.map(function(e){return'"'+e.content+'"'}).join(", ")+"]"}}}},77:function(e,t){},78:function(e,t){},79:function(e,t){},82:function(e,t,s){s(78);var n=s(14)(s(53),s(85),"data-v-4e1aa818",null);e.exports=n.exports},83:function(e,t,s){s(77);var n=s(14)(s(54),s(84),"data-v-233a0e32",null);e.exports=n.exports},84:function(e,t){e.exports={render:function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"try"},[s("el-card",[s("div",{staticClass:"clearfix messages-header",slot:"header"},[s("span",{staticClass:"header-text"},[e._v("\n        Test Queue"),s("br"),e._v(" "),s("el-tag",{attrs:{type:"primary"}},[e._v(e._s(e.testSquad))]),e._v(" "),s("el-tag",{attrs:{type:e.testingTagType}},[e._v(e._s(e.testingStateString))])],1),e._v(" "),s("el-button",{staticStyle:{float:"right"},attrs:{disabled:e.testingAvailable,type:"primary"},on:{click:e.pushMessage}},[e._v("Push Message")])],1),e._v(" "),e._l(e.messageLogs,function(t,n){return s("div",{staticClass:"text item"},[e._v("\n      "+e._s(t.time)+" "+e._s(t.messages)+"\n    ")])})],2)],1)},staticRenderFns:[]}},85:function(e,t){e.exports={render:function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",[s("div",{staticClass:"intro"},[s("el-row",[s("el-col",{attrs:{span:4}},[s("br")]),e._v(" "),s("el-col",{attrs:{span:8}},[s("strong",{staticClass:"head"},[e._v("SQS")]),e._v(" "),s("p",{staticClass:"desc"},[e._v("\n      Simple Queue Service for\n    ")]),e._v(" "),s("p",{staticClass:"desc"},[e._v("\n      Easily scaling, Ordered Delivery\n    ")]),e._v(" "),s("el-button",{staticClass:"link-button",attrs:{size:"large",type:"primary"}},[s("strong",[e._v("Get Started")])]),e._v(" "),s("el-button",{staticClass:"link-button",attrs:{size:"large",type:"primary"}},[s("strong",[e._v("Look the Design")])])],1),e._v(" "),s("el-col",{attrs:{span:8}},[s("try")],1),e._v(" "),s("el-col",{attrs:{span:4}})],1)],1),e._v(" "),s("br"),e._v(" "),s("el-row",[s("el-col",{attrs:{span:4}},[s("br")]),e._v(" "),s("el-col",{attrs:{span:16}},[s("h2",[s("strong",[e._v("Why a another Queue Service?")])]),e._v(" "),s("p",[e._v("\n      To be honest, firstly for learning the queue service.\n    ")]),e._v(" "),s("h2",[s("strong",[e._v("Design Overview")])]),e._v(" "),s("p",[e._v("\n      Master and node.\n    ")]),e._v(" "),s("section",[s("h2",[s("strong",[e._v("Get Started")])]),e._v(" "),s("p",[e._v("\n        Process overview.\n      ")]),e._v(" "),s("p",[e._v("\n        1. Get the `accessKey` and `secretKey`\n      ")]),e._v(" "),s("p",[e._v("\n        2. Apply a server node.\n      ")]),e._v(" "),s("p",[e._v("\n        3. Apply a message id.\n      ")]),e._v(" "),s("p",[e._v("\n        4. Push a message.\n      ")]),e._v(" "),s("p",[e._v("\n        5. Pull messages.\n      ")])])]),e._v(" "),s("el-col",{attrs:{span:4}},[s("br")])],1)],1)},staticRenderFns:[]}},86:function(e,t){e.exports={render:function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{attrs:{id:"app"}},[s("div",{staticClass:"header"},[s("el-row",[s("el-col",{attrs:{span:4}},[s("br")]),e._v(" "),s("el-col",{attrs:{span:10}},[s("a",{staticClass:"nav-home",attrs:{href:"/"}},[s("svg",{staticClass:"nav-logo",attrs:{width:"51px",height:"48px",viewBox:"0 0 51 48",version:"1.1",xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink"}},[s("desc",[e._v("Created with Sketch.")]),e._v(" "),s("defs"),e._v(" "),s("g",{attrs:{id:"Page-1",stroke:"none","stroke-width":"1",fill:"none","fill-rule":"evenodd"}},[s("path",{attrs:{d:"M35.3446225,37.1000297 C34.9242758,37.9487199 34.8982216,38.83043 35.1811085,39.7195914 C35.562712,40.9190616 36.3539757,41.7714769 37.4465969,42.3452029 C38.3695274,42.8298456 39.3559591,43.0060687 40.3883082,42.9105841 C41.9681286,42.7644525 43.2946136,42.0821119 44.392945,40.9358438 C45.4000336,39.8848195 46.0190148,38.6256018 46.2848439,37.1925001 C46.6399171,35.2780883 46.3352439,33.4590112 45.4350375,31.7464878 C44.4728684,29.9160625 43.0316071,28.5903998 41.1840293,27.7124007 C40.2227474,27.255594 39.2116192,26.9696045 38.1552054,26.855418 C35.4897104,26.5673289 33.0641918,27.2325094 30.8687137,28.7686824 C29.3458717,29.8341994 28.1729597,31.219009 27.3013408,32.8702751 C26.8981941,33.6339995 26.0829158,36.29869 25.9937832,36.6957462 C25.6440467,38.2539299 25.0728315,39.7202141 24.2629869,41.0916707 C23.3445843,42.6469193 22.1798545,43.9823591 20.7661779,45.0922955 C19.6192822,45.9927585 18.3627131,46.6957241 16.9951854,47.1917374 C16.016619,47.5466791 15.0122075,47.785878 13.9784797,47.9084398 C13.0866758,48.0141748 12.1942776,48.0299225 11.3026799,47.9473798 C9.12114488,47.7454735 7.10981904,47.0425063 5.29030231,45.7990674 C3.5670481,44.6214178 2.22390239,43.0974402 1.2702588,41.2267756 C0.714956941,40.1374941 0.344709211,38.983575 0.147761931,37.7736635 C-0.0267598078,36.7015213 -0.0460549272,35.6271076 0.0865273698,34.5476295 C0.302062521,32.7927279 0.890153194,31.1793937 1.87798907,29.7225543 C2.91837129,28.1882353 4.26939248,27.015095 5.93837086,26.2246254 C6.71346674,25.8575246 7.52261262,25.5980019 8.36606978,25.45372 C9.25226339,25.3021009 10.141473,25.2804674 11.0358933,25.3890935 C12.0940627,25.5176132 13.094503,25.8302801 14.0388772,26.3239916 C14.800406,26.722069 15.486295,27.2292801 16.0935916,27.8436899 C16.6451065,28.4017157 17.1101361,29.0258857 17.4801745,29.7200325 C17.8585962,30.4298969 18.1234086,31.1819814 18.2675374,31.9769396 C18.4316222,32.8817796 18.442734,33.7857446 18.2726121,34.6914309 C18.0755568,35.7405316 17.6492551,36.6815736 17.0122963,37.5325728 C16.3704874,38.3900581 15.5560695,39.0176138 14.589764,39.4475417 C14.0238404,39.6993598 13.4255181,39.8349391 12.8097311,39.8758518 C12.0077709,39.9291162 11.2343087,39.8031841 10.4897731,39.4845938 C9.81155091,39.1943857 9.23771338,38.7643588 8.76811176,38.1969182 C8.28006617,37.6071735 7.9771132,36.9267041 7.86501963,36.1622379 C7.80319738,35.7405899 7.81772033,35.3223913 7.89532084,34.903497 C8.03037462,34.1743564 8.36808592,33.5553574 8.88916655,33.0420495 C9.22083778,32.7153229 9.6110164,32.4766743 10.0526403,32.3219054 C10.4995173,32.1652922 10.9555332,32.1050269 11.4235893,32.1532905 C11.9525288,32.207826 12.424363,32.404102 12.8252655,32.7652104 C12.9594689,32.8860908 13.0825413,33.0158084 13.1858276,33.1647647 C13.2790265,33.0006518 13.2288506,32.7809061 13.1539602,32.566535 C12.8797388,31.7814797 12.3506919,31.2208095 11.6304305,30.8456388 C10.6558631,30.3380148 9.63241356,30.2813298 8.59313113,30.5922921 C7.24666427,30.9951163 6.25649672,31.8611817 5.56658843,33.0837526 C5.009793,34.0704095 4.76480615,35.1450407 4.79364379,36.2757996 C4.82501554,37.5053787 5.16568922,38.6509542 5.80110364,39.7014788 C6.76589923,41.2965807 8.14209527,42.3783373 9.87179875,43.0055771 C10.9386029,43.3924482 12.0419186,43.547057 13.1742162,43.4777973 C15.3213531,43.3465076 17.1952252,42.5395818 18.8040002,41.1070227 C20.4478226,39.6432924 21.4836091,37.8061393 21.9903066,35.6576071 C22.2284739,34.6476858 22.313435,33.6203383 22.2562417,32.585844 C22.1452016,30.5778886 21.5503974,28.7223461 20.4972676,27.0164437 C20.0722374,26.3279692 19.5755141,25.6980975 19.0185945,25.1154453 C18.9233103,25.0157283 17.3846436,23.5018008 16.8746056,22.9254808 C15.7797602,21.6883677 14.9147259,20.3046004 14.2802408,18.7728524 C13.7988458,17.6106507 13.4761951,16.4039308 13.3093682,15.1561531 C13.1641949,14.0701731 13.1427684,12.9793738 13.2553871,11.8889883 C13.486078,9.6553233 14.1967249,7.59282286 15.4285993,5.71999448 C16.4609145,4.15057411 17.7660944,2.86408934 19.3486616,1.86994365 C20.6623937,1.04470019 22.0812262,0.489289227 23.6020584,0.204144755 C24.5018856,0.0354410912 25.4103107,-0.03457068 26.3234119,0.0162356056 C28.829684,0.155683271 31.0786644,0.990320622 33.0301938,2.60514215 C34.2571235,3.6203718 35.2171829,4.85233291 35.9042004,6.29807409 C36.3771885,7.29346195 36.6840724,8.34125812 36.8245234,9.43770927 C36.9268817,10.2368717 36.9400906,11.0378621 36.8479219,11.8381669 C36.6530028,13.5308656 36.07425,15.0734334 35.071414,16.448232 C34.2946422,17.5130667 33.3324256,18.3633868 32.1835053,18.9923449 C31.344379,19.4517125 30.4506737,19.753824 29.5054449,19.8988056 C28.7937126,20.0079946 28.0806139,20.0201221 27.3680378,19.9293377 C26.4443416,19.8116672 25.5700068,19.5295255 24.7460137,19.0872386 C24.0443553,18.7106224 23.4242336,18.2271674 22.886891,17.6374023 C22.3861029,17.0877318 21.9956607,16.4631782 21.7018005,15.7769857 C21.4733955,15.243625 21.3106947,14.6910043 21.2359783,14.113905 C21.1784617,13.6696488 21.1827889,13.2236615 21.2136728,12.7774554 C21.2518926,12.2251812 21.3928351,11.6979497 21.592195,11.1851838 C21.864526,10.48472 22.2697187,9.87198274 22.8073148,9.35386135 C23.3861135,8.79602385 24.0566803,8.3997568 24.8360072,8.19173797 C25.2321913,8.08598814 25.6328171,8.03790556 26.0374551,8.04750073 C26.3877933,8.0558104 26.7364893,8.10649056 27.0733014,8.22003494 C27.3788255,8.32301488 27.6789374,8.43395544 27.9539622,8.60922866 C28.4772514,8.94275513 28.9062228,9.36424752 29.2128811,9.91380018 C29.4604456,10.3574321 29.613967,10.8284952 29.6431601,11.3349991 C29.6991483,12.3066333 29.4064988,13.1578337 28.7180714,13.844372 C28.3395894,14.2218495 27.8669959,14.4326279 27.3351726,14.4966226 C27.7585937,14.8117103 28.2415112,14.8681299 28.7329149,14.7932465 C29.9714815,14.6045139 30.8348647,13.881094 31.4155896,12.7914682 C31.9297364,11.8267599 32.0550021,10.7900046 31.8644598,9.71989656 C31.569094,8.06104129 30.6714948,6.79781413 29.3034902,5.87151635 C28.11877,5.06931988 26.7980049,4.71978065 25.3779087,4.75833698 C23.5444478,4.8080902 21.9405231,5.47079103 20.5525113,6.66855075 C19.1237512,7.9014543 18.2227751,9.46749306 17.7712267,11.3045977 C17.5627396,12.1528319 17.4791971,13.0169971 17.5163648,13.8886001 C17.5998644,15.847007 18.2009999,17.634011 19.3022157,19.2452648 C20.4236273,20.8861135 21.896163,22.1048999 23.6809882,22.9419251 C24.6485272,23.3956922 25.6621638,23.7000473 26.7151821,23.8587277 C28.0153504,24.0546527 29.3155037,24.0429792 30.6086103,23.7876646 C30.9646847,23.7173596 33.1473738,23.1086626 33.7430589,22.9611424 C34.5132552,22.7703782 35.2934693,22.6426119 36.084824,22.5801444 C36.7216029,22.5298989 37.3581805,22.5183019 37.995371,22.55413 C39.8781287,22.6599827 41.6735318,23.1147302 43.3751863,23.9417948 C45.0563596,24.7589385 46.5167176,25.8719816 47.7457419,27.2925963 C49.1809006,28.9514996 50.1583763,30.8524208 50.6620312,32.9990031 C50.9770287,34.341603 51.0711147,35.7014985 50.9473101,37.0766791 C50.7655799,39.0954643 50.1417503,40.9636683 49.0462584,42.6588086 C47.7773831,44.6221909 46.0828144,46.0880297 43.9433759,47.0101286 C42.9121447,47.4546051 41.8374937,47.7282455 40.7203843,47.835448 C39.9095256,47.9132489 39.1007565,47.9019971 38.2963002,47.7826683 C35.997754,47.4417375 34.0194302,46.4497442 32.4171523,44.7382463 C31.1651803,43.4009932 30.3674687,41.8159227 30.0539078,39.9969059 C29.7849521,38.4365685 29.9254701,36.9102616 30.5002554,35.4308427 C30.8616982,34.5004764 31.3763114,33.6684423 32.0512686,32.9376551 C32.5969852,32.3467725 33.2178805,31.8519024 33.9173409,31.4640393 C34.6056806,31.0823342 35.3384357,30.8220046 36.1166179,30.6959403 C37.1161922,30.5340324 38.1004395,30.6077484 39.0675741,30.896201 C39.5599469,31.0430722 40.0182067,31.2737413 40.4566274,31.5464411 C41.164696,31.9868423 41.7341784,32.5688837 42.1941392,33.2640766 C42.2931929,33.4137728 42.3794497,33.5735079 42.4584665,33.735215 C42.8595608,34.5560773 43.0565155,35.4181029 42.9876415,36.3399894 C42.9339932,37.0580587 42.7276649,37.7240241 42.3588417,38.3346886 C42.1695528,38.6480723 41.943902,38.9382016 41.6640005,39.1795957 C41.4600237,39.3554918 41.2578001,39.5302964 41.024136,39.6691521 C40.3322122,40.0803577 39.5872224,40.2465476 38.7949884,40.1546632 C38.1420359,40.0789088 37.5579017,39.8219369 37.0568528,39.3807408 C36.7193089,39.0835368 36.4523204,38.7346074 36.2602778,38.3261189 C35.944409,37.6542921 35.9093728,36.9722755 36.1857043,36.2783094 C36.2073287,36.2239994 36.2427924,36.1738626 36.2412185,36.0888366 C35.8324126,36.3458937 35.5485144,36.6884365 35.3446225,37.1000297 Z",id:"Fill-3",fill:"#3CDBFF"}})])])])]),e._v(" "),s("el-col",{attrs:{span:3}},[s("br")]),e._v(" "),s("el-col",{attrs:{span:3}},[s("div",{staticClass:"header-right"},[s("a",{staticClass:"link",attrs:{href:"https://github.com/go-sqs",type:"text"}},[e._v("Github")]),e._v(" "),s("a",{staticClass:"link",attrs:{href:"https://github.com/go-sqs",type:"text"}},[e._v("V0.1.0")])])]),e._v(" "),s("el-col",{attrs:{span:4}},[s("br")])],1)],1),e._v(" "),s("router-view",{staticClass:"router"})],1)},staticRenderFns:[]}},89:function(e,t){}},[51]);
//# sourceMappingURL=app.a104fbfda3536ce26802.js.map