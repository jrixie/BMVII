<template>
  <div id="app">
    <input v-model="cname">
    
    <!-- Have default option -->
    <select v-model="ccond">
      <option disabled value="">Select lowest preferred condition</option>
      <option>Near Mint</option>
      <option>Lightly Played</option>
      <option>Moderately Played</option>
      <option>Heavily Played</option>
      <option>Damaged</option>
    </select>

    <money v-model="price" v-bind="money"></money> 
  </div>
</template>

<script>

import {Money} from 'v-money'

export default {
  name: 'app',
  components: {
    Money
  },
  data() {
    return {
      price: 0.00,
        money: {
          decimal: '.',
          thousands: ',',
          prefix: '$ ',
          suffix: '',
          precision: 2,
          masked: false
        }
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
