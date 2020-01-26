<template>
  <div id="app">
    <p v-if="errors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="error in errors" v-bind:key="error.message">{{ error.message }}</li>
      </ul>
    </p>

    <input v-model="email" placeholder="email">

    <input v-model="cname" placeholder="cardname">
    
    <!-- Have default option -->
    <select v-model="ccond">
      <option disabled value="">Select lowest preferred condition</option>
      <option>Near Mint</option>
      <option>Lightly Played</option>
      <option>Moderately Played</option>
      <option>Heavily Played</option>
      <option>Damaged</option>
    </select>

    <select v-model="ineq">
      <option v-for="option in ineqOptions" v-bind:key="option.value">
        {{ option.text }}
      </option>
    </select>

    <money v-model="setprice" v-bind="money"></money> 

    <button v-on:click="formSubmit">Submit</button>
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
      post: {},
      price: 0.00,
      money: {
        decimal: '.',
        thousands: ',',
        prefix: '$ ',
        suffix: '',
        precision: 2,
        masked: false
      },
      email: '',
      cname: '',
      ccond: '',
      ineq: '',
      setprice: '',
      ineqOptions: [
        { text: '>', value: '0' },
        { text: '<', value: '1' },
      ],
      errors: [],
      reg: /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@(([[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,24}))$/
    }
  },
  methods: {
            formSubmit(e) {
              e.preventDefault();
              this.errors = [];
              if (!this.validEmail(this.email)) {
                this.errors.push({message: 'Valid email required.'});
              } else {
                let currentObj = this;
                this.post = {
                    email: this.email,
                    cardName: this.cname,
                    cardCondition: this.ccond,
                    priceCondition: this.ineq,
                    priceThreshold: this.setprice
                };
                this.axios.post('http://localhost:4000/api/alert/create/', this.post)
                .then(function (response) {
                    currentObj.output = response.data;
                })
                .catch(function (error) {
                    currentObj.output = error;
                });
              }
            },
            validEmail(email) {
              return this.reg.test(email);
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
