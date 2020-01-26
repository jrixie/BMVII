<template>
  <div id="app">
    <h1>
      MTG Price Alerts
    </h1>

    <p v-if="!validEmail">
      <b>Please enter a valid email</b>
    </p>

    <div class = "fields">
      <input v-model="email" placeholder="Email" class="entry">

      <input v-model="cname" placeholder="Card Name" class="entry">
      

      <!-- Have default option -->
      <select v-model="ccond" class="entry">
        <option disabled value="">Select Lowest Preferred Condition</option>
        <option>Near Mint</option>
        <option>Lightly Played</option>
        <option>Moderately Played</option>
        <option>Heavily Played</option>
        <option>Damaged</option>
      </select>



      <select v-model="ineq" class="entry">
        <option disabled value="">Greater Than or Less Than</option>
        <option v-for="option in ineqOptions" v-bind:key="option.value" v-bind:value="option.value">
          {{ option.text }}
        </option>
      </select>


      <money v-model="setprice" v-bind="money" class="entry"></money>

      <div><button v-on:click="formSubmit" style="width:100px;margin:10px;">Submit</button></div>

      <img src="./assets/mtglogo.png" style="padding-top:50px;">
    </div>
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
        { text: '>', value: 0 },
        { text: '<', value: 1 },
      ],
      validEmail: true,
      reg: /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@(([[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,24}))$/
    }
  },
  methods: {
            formSubmit(e) {
              e.preventDefault();
              this.validEmail = true;
              if (!this.checkEmail(this.email)) {
                this.validEmail = false;
              } else {
                let currentObj = this;
                this.post = {
                    email: this.email,
                    cardName: this.cname,
                    cardCondition: this.ccond,
                    priceCondition: this.ineq,
                    priceThreshold: this.setprice
                };
                this.axios.post('http://mtgdrop.tech:4000' + '/api/alert/create/', this.post)
                .then(function (response) {
                    currentObj.output = response.data;
                })
                .catch(function (error) {
                    currentObj.output = error;
                });
              }
            },
            checkEmail(email) {
              return this.reg.test(email);
           }
      }
}
</script>

<style lang="scss">
@font-face {
    font-family: 'Beleren';
    src: url(assets/Beleren2016-Bold.ttf);
}

$fontColor: white;
$backgroundColor: #384048 ;
$borderSize: 200px;

* {
      margin: 0;
  }

#app {
  height: 1000px;
  font-family: Beleren, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: $fontColor;
  background-color: $backgroundColor;
  //#2c3e50;

  .fields {
    text-align:center;
    margin: auto;
    border-radius: 25px;
    background: black;
    padding: 20px;
    width: 35%;
    height: 250px;
  }

  .entry {
    width: 300px;
    margin:10px;
  }
}
</style>
