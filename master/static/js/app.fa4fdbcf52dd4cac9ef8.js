webpackJsonp([1],{23:function(e,t,s){"use strict";var n=s(2),a=s(83),o=s(78),i=s.n(o);n.default.use(a.a),t.a=new a.a({routes:[{path:"/",name:"Intro",component:i.a}]})},25:function(e,t){},26:function(e,t,s){s(77);var n=s(14)(s(52),s(82),null,null);e.exports=n.exports},51:function(e,t,s){"use strict";Object.defineProperty(t,"__esModule",{value:!0});var n=s(2),a=s(24),o=s.n(a),i=s(25),r=(s.n(i),s(26)),c=s.n(r),l=s(23),d=s(27),u=s(28),g=s.n(u);n.default.use(g.a),n.default.use(o.a),n.default.use(d.a),n.default.config.productionTip=!1,new n.default({el:"#app",router:l.a,template:"<App/>",components:{App:c.a}})},52:function(e,t,s){"use strict";Object.defineProperty(t,"__esModule",{value:!0}),t.default={name:"app"}},53:function(e,t,s){"use strict";Object.defineProperty(t,"__esModule",{value:!0});var n=s(79),a=s.n(n);t.default={name:"intro",components:{Try:a.a},data:function(){var e="";return{masterAddr:e,testAccessKey:"test",testSecretKey:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2dpdGh1Yl9sb2dpbl9rZXkiOiJ0ZXN0In0.Rmg9qFF_9YwBpIYD_Utd4L89bbIP7Vl9yyzehdos1L8"}},methods:{sendKeysToStagzar:function(){var e=this;this.$prompt("Then input your github.com login","Star the github.com/Focinfi/gosqs",{confirmButtonText:"Send Keys",cancelButtonText:"Cancel",inputValidator:function(e){return null!==e&&void 0!==e&&""!==e||"Empty login"},beforeClose:function(t,s,n){if("confirm"!==t)return void n();s.confirmButtonLoading=!0,e.$http.get(e.masterAddr+"/sendGithubEmailSecretKey/"+s.inputValue).then(function(t){var a=t.body;if(1e3===a.code)return n(),setTimeout(function(){s.confirmButtonLoading=!1},300),void e.$notify({type:"success",title:"Notification",message:"Sent to "+a.data.email,duration:0});e.$message({type:"warning",message:'"'+s.inputValue+'" is '+a.message}),s.confirmButtonLoading=!1},function(t){n(),s.confirmButtonLoading=!1,e.$message({type:"info",message:"Failed to send the keys"})})}})}}}},54:function(e,t,s){"use strict";Object.defineProperty(t,"__esModule",{value:!0});var n=s(55),a=s.n(n);t.default={name:"try",props:{masterAddr:String,accessKey:String,secretKey:String},data:function(){return{testSquad:"gosqs.org",newMessage:"",messageLogs:[],servingNode:"",token:"",nextId:0}},created:function(){this.applyNode()},mounted:function(){console.log("mounted"),setInterval(this.pullMessage,1e3)},computed:{apiJSONparams:function(){return a()({queue_name:this.accessKey,squad_name:this.testSquad,token:this.token,message_id:this.nextId,content:this.newMessage,size:1})},testingTagType:function(){return""===this.servingNode?"gray":"success"},testingUnavailable:function(){return""===this.servingNode},testingStateString:function(){return""===this.servingNode?"unavailabe":"availabe"}},methods:{applyNode:function(){var e=this,t=a()({access_key:this.accessKey,secret_key:this.secretKey,queue_name:this.accessKey,squad_name:this.testSquad});this.$http.post(this.masterAddr+"/applyNode",t).then(function(t){var s=t.body;if(1e3===s.code){var n=s.data;e.servingNode=n.node,e.token=n.token}})},pullMessage:function(){var e=this;this.testingUnavailable||this.$http.post(this.servingNode+"/messages",this.apiJSONparams).then(function(t){var s=t.body.data,n=e.timeFormat(new Date),a={time:n,messages:[{message_id:0}]};null!==s&&s.length>0&&(a.messages=s,e.reportReceived()),e.messageLogs.push(a),e.messageLogs.length>4&&e.messageLogs.shift()},function(t){e.servingNode="",console.log(t)})},pushMessage:function(){var e=this,t=this;this.$prompt("Input Message","Ready To Push",{confirmButtonText:"Push it",cancelButtonText:"Cancel",inputValidator:function(e){return null===e||void 0===e||""===e?"Can not push empty message":e.length>20?"Content is too long":(console.log(e),t.newMessage=e,!0)}}).then(function(s){s.value;e.applyMessageID(function(){t.$http.post(t.servingNode+"/message",t.apiJSONparams).then(function(e){t.newMessage="",console.log(e.data),t.$message({type:"success",message:"Message pushed"})},function(e){t.$message({type:"warning",message:"Failed to push message"})})})})},applyMessageID:function(e){var t=this;this.$http.post(this.servingNode+"/messageID",this.apiJSONparams).then(function(s){console.log(s.data);var n=s.body.data;t.nextId=n.message_id_end,console.log("nextID",t.nextId),e()},function(e){})},reportReceived:function(){this.$http.post(this.servingNode+"/receivedMessageID",this.apiJSONparams).then(function(e){},function(e){console.log(e)})},messagesLogsSlice:function(e){console.log(this.messageLogs.length)},timeFormat:function(e){return"["+e.toLocaleTimeString("en-GB")+"]"},messagesEntryFormat:function(e){return"["+e.map(function(e){return 0===e.message_id?"wait for new messages":e.message_id+'-"'+e.content+'"'}).join(", ")+"]"}}}},75:function(e,t){},76:function(e,t){},77:function(e,t){},78:function(e,t,s){s(76);var n=s(14)(s(53),s(81),"data-v-4e1aa818",null);e.exports=n.exports},79:function(e,t,s){s(75);var n=s(14)(s(54),s(80),"data-v-233a0e32",null);e.exports=n.exports},80:function(e,t){e.exports={render:function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"try"},[s("el-card",[s("div",{staticClass:"clearfix messages-header",slot:"header"},[s("span",{staticClass:"header-text"},[e._v("\n        Test Queue"),s("br"),e._v(" "),s("el-tag",{attrs:{type:"primary"}},[e._v(e._s(e.testSquad))]),e._v(" "),s("el-tag",{attrs:{type:e.testingTagType}},[e._v(e._s(e.testingStateString))])],1),e._v(" "),s("el-button",{staticStyle:{float:"right"},attrs:{disabled:e.testingUnavailable,type:"primary"},on:{click:e.pushMessage}},[e._v("Push Message")])],1),e._v(" "),e._l(e.messageLogs,function(t,n){return s("div",{staticClass:"text item"},[e._v("\n      "+e._s(t.time)+" "+e._s(e.messagesEntryFormat(t.messages))+"\n    ")])})],2)],1)},staticRenderFns:[]}},81:function(e,t){e.exports={render:function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",[s("div",{staticClass:"intro"},[s("el-row",[s("el-col",{attrs:{span:4}},[s("br")]),e._v(" "),s("el-col",{attrs:{span:8}},[s("strong",{staticClass:"head"},[e._v("GoSQS")]),e._v(" "),s("p",{staticClass:"desc"},[e._v("\n      Simple Queue Service for\n    ")]),e._v(" "),s("p",{staticClass:"desc"},[e._v("\n      Scalability and Ordered Delivery\n    ")]),e._v(" "),s("el-button",{directives:[{name:"scroll-to",rawName:"v-scroll-to",value:{el:"#get-started",offset:-80},expression:"{el: '#get-started', offset: -80}"}],staticClass:"link-button",attrs:{size:"large",type:"primary"}},[s("strong",[e._v("Get Started")])]),e._v(" "),s("el-button",{directives:[{name:"scroll-to",rawName:"v-scroll-to",value:{el:"#design-overview",offset:-80},expression:"{el: '#design-overview', offset: -80}"}],staticClass:"link-button",attrs:{size:"large",type:"primary"}},[s("strong",[e._v("Look the Design")])])],1),e._v(" "),s("el-col",{attrs:{span:8}},[s("try",{attrs:{masterAddr:e.masterAddr,accessKey:e.testAccessKey,secretKey:e.testSecretKey}})],1),e._v(" "),s("el-col",{attrs:{span:4}})],1)],1),e._v(" "),s("br"),e._v(" "),s("el-row",[s("el-col",{attrs:{span:4}},[s("br")]),e._v(" "),s("el-col",{attrs:{span:16}},[s("h2",[s("strong",[e._v("Why another Queue Service?")])]),e._v(" "),s("p",[e._v("\n      While playing the golang programming language in distributed system.\n      A queue service project called gosqs with focus on  \n      "),s("strong",[e._v("reliability")]),e._v(",\n      "),s("strong",[e._v("scalability")]),e._v(" and\n      "),s("strong",[e._v("ordered delivery")]),e._v(" was born.\n      "),s("br"),e._v(" "),s("br"),e._v("\n      Compare with other message queue implementations, the currency performance is not the priority of gosqs. Furthermore, gosqs is designed for the more strict scenarios. \n    ")]),e._v(" "),s("h2",{attrs:{id:"design-overview"}},[s("strong",[e._v("Design Overview")])]),e._v(" "),s("p",[e._v("\n    gosqs uses the\n      "),s("a",{attrs:{href:"https://github.com/coreos/etcd",target:"_blank"}},[e._v("etcd")]),e._v("\n      to store the meta data of cluster and queue, like the current available node servers.\n    "),s("br"),e._v(" "),s("br"),e._v("\n    To decouple the dependency of underlying storage of message, gosqs uses the \n      "),s("a",{attrs:{href:"https://github.com/Focinfi/oncekv",target:"_blank"}},[e._v("oncekv")]),e._v("\n      which is a combination of groupcahe and raft-boltdb.\n    ")]),e._v(" "),s("img",{staticClass:"architeture",attrs:{src:"http://on78mzb4g.bkt.clouddn.com/gosqs_architeture.png",alt:"architeture"}}),e._v(" "),s("br"),e._v(" "),s("p",[e._v("\n    Every message in a queue has its own sequence id, every queue has a number of squads, which is a record for the last processed message id in its coresponding queue.  \n    ")]),e._v(" "),s("img",{staticClass:"queue_and_squad",attrs:{src:"http://on78mzb4g.bkt.clouddn.com/gosqs_queue_squad.png",alt:"architeture"}}),e._v(" "),s("section",[s("h2",{attrs:{id:"get-started"}},[s("strong",[e._v("Get Started")])]),e._v(" "),s("p",[e._v("\n        1. Perpare the keys: \n        ")]),s("div",{staticClass:"indent"},[e._v("\n          Access Key:\n          "),s("code",[e._v('"'+e._s(e.testAccessKey)+'"')]),e._v(" "),s("br"),e._v(" "),s("br"),e._v("\n          Secret Key:\n          "),s("el-input",{staticClass:"secret-key-input",attrs:{size:"small"},nativeOn:{click:function(e){e.target.select()}},model:{value:e.testSecretKey,callback:function(t){e.testSecretKey=t},expression:"testSecretKey"}}),e._v(" "),s("br"),e._v(" "),s("br"),e._v("\n          Or, \n          "),s("el-button",{attrs:{type:"text"},on:{click:e.sendKeysToStagzar}},[e._v("Get your own keys")])],1),e._v(" "),s("p"),e._v(" "),s("p",[e._v("\n        2. Apply a server node:\n        "),s("pre",{staticClass:"code indent"},[e._v("\n  "),s("code",[e._v('POST http://master.gosqs.org/applyNode\n  Body:\n    {\n      "access_key": "test", \n      "secret_key": "secret.key.in.step.1",\n      "queue_name": "foo",\n      "squad_name": "bar"\n    }')]),e._v("\n        ")])]),s("p",{staticClass:"indent"},[e._v("\n          You will get the response body like:\n          "),s("pre",{staticClass:"code"},[e._v("\n  "),s("code",[e._v('{\n      "code": 1000,\n      "message": "",\n      "data": {\n          "node": "http://node-1.gosqs.org",\n          "token": "token.token.token"\n      }\n  }')]),e._v("\n          ")]),e._v('\n          The "token" will be used in the subsequent requests.\n        ')]),e._v(" "),s("p"),e._v(" "),s("p",[e._v("\n        3. Apply a message id:\n        ")]),s("p",{staticClass:"indent"},[e._v('\n          Every new message need a sequence id.\n          Specify the "size" parameter to determine the range of id.\n          '),s("pre",{staticClass:"code"},[e._v("\n  "),s("code",[e._v('POST http://node-1.gosqs.org/messageID\n  Body:\n    {\n      "queue_name": "foo",\n      "token": "token.from.step.2",\n      "size": 1\n    }')]),e._v("\n          ")]),e._v("\n          Response as following which indicates the id is 5.\n          "),s("pre",{staticClass:"code"},[e._v("\n  "),s("code",[e._v('{\n    "code": 1000,\n    "message": "",\n    "data": {\n        "message_id_begin": 5,\n        "message_id_end": 5\n    }\n  }')]),e._v("\n          ")])]),e._v(" "),s("p"),e._v(" "),s("p",[e._v("\n        4. Push a message:\n        ")]),s("p",{staticClass:"indent"},[s("pre",{staticClass:"code"},[e._v("\n  "),s("code",[e._v('POST http://node-1.gosqs.org/message\n  Body:\n    {\n      "queue_name": "foo",\n      "token": "token.from.step.2",\n      "message_id": 5,\n      "content": "hello gosqs"\n    }')]),e._v("\n          ")]),e._v('\n          Check the value of "code" in the response body to detect if the message has been pushed, 1000 is ok.\n          '),s("pre",{staticClass:"code"},[e._v("\n  "),s("code",[e._v('{\n    "code": 1000,\n    "message": "",\n    "data": null\n  }')]),e._v("\n          ")])]),e._v(" "),s("p"),e._v(" "),s("p",[e._v("\n        5. Pull messages:\n        ")]),s("p",{staticClass:"indent"},[s("pre",{staticClass:"code"},[e._v("\n  "),s("code",[e._v('POST http://node-1.gosqs.org/messages\n  Body:\n    {\n      "queue_name": "foo",\n      "squad_name": "bar",\n      "token": "token.from.step.2",\n    }')]),e._v("\n          ")]),e._v('\n          Messages will be in the "data" parameter.\n          '),s("pre",{staticClass:"code"},[e._v("\n  "),s("code",[e._v('{\n    "code": 1000,\n    "message": "",\n    "data": [\n      {"message_id": 4, "content": "hello from china"},\n      {"message_id": 5, "content": "hello gosqs"}\n    ]\n  }')]),e._v("\n          ")])]),e._v(" "),s("p")])]),e._v(" "),s("el-col",{attrs:{span:4}},[s("br")])],1)],1)},staticRenderFns:[]}},82:function(e,t){e.exports={render:function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{attrs:{id:"app"}},[s("div",{staticClass:"header"},[s("el-row",[s("el-col",{attrs:{span:4}},[s("br")]),e._v(" "),s("el-col",{attrs:{span:10}},[s("a",{staticClass:"nav-home",attrs:{href:"/"}},[s("img",{staticClass:"nav-logo",attrs:{src:"http://on78mzb4g.bkt.clouddn.com/gosqs.org.logo.svg"}})])]),e._v(" "),s("el-col",{attrs:{span:3}},[s("br")]),e._v(" "),s("el-col",{attrs:{span:3}},[s("div",{staticClass:"header-right"},[s("a",{staticClass:"link",attrs:{href:"https://github.com/Focinfi/gosqs",type:"text"}},[e._v("Github")]),e._v(" "),s("a",{staticClass:"link",attrs:{href:"https://github.com/Focinfi/gosqs",type:"text"}},[e._v("V0.1.0")])])]),e._v(" "),s("el-col",{attrs:{span:4}},[s("br")])],1)],1),e._v(" "),s("router-view",{staticClass:"router"})],1)},staticRenderFns:[]}},86:function(e,t){}},[51]);
//# sourceMappingURL=app.fa4fdbcf52dd4cac9ef8.js.map