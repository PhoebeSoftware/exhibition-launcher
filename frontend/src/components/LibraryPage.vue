<template>
  <div class="page">
    <h1>Library</h1>
    <p>Your saved content appears here.</p>
    <div>
      <div v-for="item in items" :key="item.id">{{ item.name }}</div>
    </div>
  </div>
</template>

<script>
import { GetAllGames } from '../../bindings/derpy-launcher072/library/library';



export default {
  name: 'LibraryPage',
  data() {
    return {
      items: []
    };
  },
  methods: {
    addItem(item) {
      this.items.push(item);
    }
  },
  async mounted() {
    const idk = await GetAllGames();

    console.log(idk);
  }
};
</script>

<style scoped>
.page {
  padding: 20px;
}

h1 {
  margin-bottom: 15px;
}
</style>