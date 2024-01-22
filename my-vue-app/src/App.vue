<template>
  <div id="app" class="container">
    <div class="header">
      <h1>VUR v.1.0.1</h1>
    </div>
    <div class="input-group">
      <input v-model="userText" placeholder="imput text" class="input-text">
      <button @click="saveText" class="btn-encrypt">Encrypt</button>
    </div>
    <div v-if="link" class="share-link">
      Share to decrypt: <a :href="link" target="_blank" class="link">{{ link }}</a>
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
    // async saveText() {
    //   try {
    //     const response = await fetch('http://localhost:8080/save', {
    //       method: 'POST',
    //       headers: {
    //         'Content-Type': 'application/json',
    //       },
    //       body: JSON.stringify({ text: this.userText }),
    //     });

    //     if (!response.ok) {
    //       throw new Error(`HTTP error! status: ${response.status}`);
    //     }

    //     const data = await response.json(); // Предполагается, что сервер возвращает JSON объект с полем uuid
    //     this.link = `http://localhost:8080/text/${data.uuid}`; // Формирование ссылки

    //     console.log('Success:', data);
    //     alert('Text saved successfully');
    //   } catch (error) {
    //     console.error('Error:', error);
    //     alert('Failed to save text');
    //   }
    // }
    async saveText() {
  try {
    const response = await fetch('http://65.109.8.6/save', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ text: this.userText }),
    });

    if (!response.ok) {
      // Выводим подробности об ошибке для дальнейшего анализа
      console.error(`HTTP error! status: ${response.status}, statusText: ${response.statusText}`);
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    if (!data.uuid) {
      throw new Error('No UUID in response');
    }

    this.link = `http://65.109.8.6/text/${data.uuid}`;
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
/* Стили для контейнера приложения */
.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background: linear-gradient(to bottom, #ADD8E6, #FFFFFF);
}

/* Стили для заголовка */
.header {
  margin-bottom: 20px;
}

/* Стили для группы ввода */
.input-group {
  display: flex;
  flex-direction: column;
  align-items: center;
}

/* Стили для поля ввода */
.input-text {
  border: 1px solid #000;
  border-radius: 5px;
  padding: 10px;
  margin-bottom: 10px;
  width: 300px; /* Можно изменить по предпочтениям */
}

/* Стили для кнопки */
.btn-encrypt {
  background-color: #4CAF50; /* Зеленый */
  border: none;
  color: white;
  padding: 15px 32px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 18px;
  margin: 4px 2px;
  cursor: pointer;
  border-radius: 5px;
}

/* Стили для ссылки */
.link {
  color: #0000EE;
  text-decoration: none;
}

/* Стили для элемента отображения ссылки */
.share-link {
  margin-top: 20px;
}
</style>
