<template>
  <div class="hello">
    <el-card>
      <div slot="header" class="clearfix messages-header">
        <span class="header-text">Test Queue</span>
        <el-button style="float: right;" type="primary" @click="pushMessage">Push Message</el-button>
      </div>
      <div v-for="o, i in messageLogs" class="text item">
        {{'Message-'+ (i + 1) +': ' + o }}
      </div>
    </el-card> 
  </div>
</template>

<script>
function httpAsync (method, theUrl, callback) {
  var xmlHttp = new XMLHttpRequest()
  xmlHttp.onreadystatechange = function () {
    if (xmlHttp.readyState === 4 && xmlHttp.status === 200) {
      callback(xmlHttp.responseText)
    }
  }
  xmlHttp.open(method, theUrl, true)
  xmlHttp.send(null)
}
// var accessKey = 'login_name'
// var secretKey = 'secret_key'

export default {
  name: 'hello',
  data () {
    return {
      testQueue: 'test',
      newMessage: '',
      messageLogs: [],
      masterAddr: 'http://127.0.0.1:5446'
    }
  },
  created () {
    console.log('created')
    this.pullMessage()
  },
  methods: {
    echo () {
      console.log('hehe')
      this.newMessage = ''
    },
    pullMessage () {
      httpAsync('POST', this.masterAddr + '/applyNode', function (responseText) {
        console.log(responseText)
      })
    },
    pushMessage () {
      this.$prompt('Input Message', 'Ready To Push', {
        confirmButtonText: 'Push it',
        cancelButtonText: 'Cancel',
        inputValidator: function (content) {
          if (content === null) {
            return 'Can not push empty message'
          } else if (content.length > 20) {
            return 'Content is too long'
          }
          return true
        }
      }).then(({ value }) => {
        this.messageLogs.push(value)
        this.$message({
          type: 'success',
          message: 'Message pushed'
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: 'Cancel Push'
        })
      })
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

.hello {
  padding: 0 35%;
  margin: 0 auto;
}

.messages-header {
  text-align: left;
}

.header-text {
  font-size: 1.5em;
  line-height: 36px;
}

.item {
  text-align: left;
}
</style>
