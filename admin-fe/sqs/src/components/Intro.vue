<template>
  <div>
  <div class="intro">
  <el-row>
    <el-col :span="4">
      <br>
    </el-col>
    <el-col :span="8">
      <strong class="head">SQS</strong>
      <p class="desc">
        Simple Queue Service for
      </p>
      <p class="desc">
        Easily scaling, Ordered Delivery
      </p>
      <el-button class="link-button" size="large" type="primary"><strong>Get Started</strong></el-button>
      <el-button class="link-button" size="large" type="primary"><strong>Look the Design</strong></el-button>
    </el-col>
    <el-col :span="8">
      <try></try>
    </el-col>
    <el-col :span="4">
    </el-col>
  </el-row>
  </div>
  <br>
  <el-row>
    <el-col :span="4">
      <br>
    </el-col>
    <el-col :span="16">
      <h2><strong>Why a another Queue Service?</strong></h2>
      <p>
        To be honest, firstly for learning the queue service.
      </p>
      <h2><strong>Design Overview</strong></h2>
      <p>
        Master and node.
      </p>
      <section>
        <h2><strong>Get Started</strong></h2>
        <p>
          Process overview.
        </p>
        <p>
          1. Get the `accessKey` and `secretKey`
        </p>
        <p>
          2. Apply a server node.
        </p>
        <p>
          3. Apply a message id.
        </p>        
        <p>
          4. Push a message.
        </p>
        <p>
          5. Pull messages.
        </p>
      </section>
    </el-col>      
    <el-col :span="4">
      <br>
    </el-col>    
  </el-row>
  </div>
</template>

<script>
import Try from './Try.vue'

export default {
  name: 'intro',
  components: {
    Try
  },
  data () {
    return {
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

.head {
  color: rgb(97, 218, 251);
  font-size: 3.5em;
  text-shadow: 0 1px 3px rgba(0,0,0,0.3);
}

.intro {
  width: 100%;
  background: rgb(45, 45, 45);
  padding: 50px 0;
}

.desc {
  font-size: 1.2em;
  color: #e9e9e9;
}

.link-button {
  background-color: rgb(204, 122, 111);
  border-color: rgb(204, 122, 111);
  font-size: 1.2em;
  text-shadow: 0 1px 3px rgba(0,0,0,0.3);
}
</style>
