<template>
  <div id="app">
    <input v-model="userText" placeholder="Введите текст">
    <button @click="saveText">Save</button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      userText: ''
    }
  },
  methods: {
    async saveText() {
      try {
        const response = await fetch('http://localhost:8080/save', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ text: this.userText }),
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.text();
        console.log('Success:', data);
        alert('Text saved');
      } catch (error) {
        console.error('Error:', error);
        alert('Failed to save text');
      }
    }
  }
}
</script>
