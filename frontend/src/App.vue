<template>
  <v-app>
    <Header :keycloak="keycloak"/>

    <v-main>
      <v-container
          class="d-flex justify-center mb-6"
          fluid>
        <!--        <v-container>-->
        <router-view/>
        <!--        </v-container>-->
      </v-container>
    </v-main>

    <Footer/>
  </v-app>
</template>

<script>
import Header from './components/template/Header';
import Footer from './components/template/Footer';

export default {
  name: 'App',

  components: {
    Header,
    Footer
  },

  props: {
    keycloak: Object
  },
  data: () => ({
    //
  }),
};
</script>

<style>
.theme--light.v-application {
  background-image: linear-gradient(#e6f0f7, #e6f0f0) !important;
}
</style>
