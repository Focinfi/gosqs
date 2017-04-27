<template>
  <div class="try">
    <el-card>
      <div slot="header" class="clearfix messages-header">
        <span class="header-text">
          Test Queue<br>
          <el-tag type="primary">{{testSquad}}</el-tag>
          <el-tag :type="testingTagType">{{testingStateString}}</el-tag>
        </span>
        <el-button :disabled="testingAvailable" style="float: right;" type="primary" @click="pushMessage">Push Message</el-button>
      </div>
      <div v-for="o, i in messageLogs" class="text item">
        {{ o.time }} {{ o.messages }}
      </div>
    </el-card> 
  </div>
</template>

<script>
let accessKey = 'Focinfi'
let secretKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyIiOiJGb2NpbmZpIn0.JRfcSY6_syXNKypblVpdD7oNmHNFDKVPVjkbNTlzEKQ'

export default {
  name: 'try',
  data () {
    return {
      testQueue: 'test',
      testSquad: 'sqsadmin',
      newMessage: '',
      messageLogs: [],
      masterAddr: process.env.SQS_ADMIN_ADDR,
      servingNode: '',
      token: '',
      nextId: 0
    }
  },
  created () {
    this.applyNode()
  },
  mounted () {
    console.log('mounted')
    setInterval(this.pullMessage, 1000)
  },
  computed: {
    apiJSONparams () {
      return JSON.stringify({
        'queue_name': this.testQueue,
        'squad_name': this.testSquad,
        'token': this.token,
        'message_id': this.nextId,
        'content': this.newMessage,
        'size': 1
      })
    },
    testingTagType () {
      return this.servingNode === '' ? 'gray' : 'success'
    },
    testingAvailable () {
      return this.servingNode === ''
    },
    testingStateString () {
      if (this.servingNode === '') {
        return 'unavailabe'
      }

      return 'availabe'
    }
  },
  methods: {
    applyNode () {
      let body = JSON.stringify({
        'access_key': accessKey,
        'secret_key': secretKey,
        'queue_name': this.testQueue,
        'squad_name': this.testSquad
      })
      this.$http.post(this.masterAddr + '/applyNode', body).then(response => {
        let body = response.body
        if (body.code !== 1000) {
          return
        }
        let data = body.data
        this.servingNode = data.node
        this.token = data.token
      }, response => {
      })
    },
    pullMessage () {
      if (this.servingNode === '') {
        this.messageLogs = [this.timeFormat(new Date()) + ' Failed to connect the node']
        return
      }
      this.$http.post(this.servingNode + '/messages', this.apiJSONparams).then(response => {
        let data = response.body.data
        // console.log(data)
        let time = this.timeFormat(new Date())
        let entry = {time: time, messages: ['No More Messages']}
        if (data.messages.length > 0) {
          entry.messages = data.messages
          this.reportReceived()
        }

        this.messageLogs.push(entry)
        if (this.messageLogs.length > 5) {
          this.messageLogs.shift()
        }
      }, response => {
        console.log(response)
      })
    },
    pushMessage () {
      let _this = this
      this.$prompt('Input Message', 'Ready To Push', {
        confirmButtonText: 'Push it',
        cancelButtonText: 'Cancel',
        inputValidator: function (content) {
          if (content === null || content === undefined || content === '') {
            return 'Can not push empty message'
          } else if (content.length > 20) {
            return 'Content is too long'
          }
          console.log(content)
          _this.newMessage = content
          return true
        }
      }).then(({ value }) => {
        // let _this = this
        this.applyMessageID(function () {
          _this.$http.post(_this.servingNode + '/message', _this.apiJSONparams).then(response => {
            _this.newMessage = ''
            console.log(response.data)
            _this.$message({
              type: 'success',
              message: 'Message pushed'
            })
          }, response => {
            _this.$message({
              type: 'fail',
              message: 'Failed to push message'
            })
          })
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: 'Cancel Push'
        })
      })
    },
    applyMessageID (push) {
      this.$http.post(this.servingNode + '/messageID', this.apiJSONparams).then(response => {
        console.log(response.data)
        let data = response.body.data
        this.nextId = data['message_id_end']
        console.log('nextID', this.nextId)
        push()
      }, response => {
      })
    },
    reportReceived () {
      this.$http.post(this.servingNode + '/receivedMessageID', this.apiJSONparams).then(response => {
      }, response => {
        console.log(response)
      })
    },
    messagesLogsSlice (n) {
      console.log(this.messageLogs.length)
    },
    timeFormat (myDate) {
      return '[' + myDate.toLocaleTimeString('en-GB') + ']'
    },
    messagesEntryFormat (messages) {
      return '[' + messages.map(function (obj) { return '"' + obj.content + '"' }).join(', ') + ']'
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}

.try {
  margin: 0 auto;
}

.messages-header {
  text-align: left;
}

.header-text {
  font-size: 1.5em;
  line-height: 34px;
}

.item {
  text-align: left;
}
</style>
