<template>
  <div id="app">
    <input v-model="userText" placeholder="Введите текст">
    <button @click="saveText">Save</button>
    <!-- Элемент для отображения ссылки -->
    <div v-if="link">
      Поделитесь этой ссылкой: <a :href="link" target="_blank">{{ link }}</a>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      userText: '',
      link: '' // Ссылка, которая будет отображаться после сохранения текста
    };
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

        const data = await response.json(); // Предполагается, что сервер возвращает JSON объект с полем uuid
        this.link = `http://localhost:8080/text/${data.uuid}`; // Формирование ссылки

        console.log('Success:', data);
        alert('Text saved successfully');
      } catch (error) {
        console.error('Error:', error);
        alert('Failed to save text');
      }
    }
  }
}
</script>

<style>
/* Ваши стили */
</style>
