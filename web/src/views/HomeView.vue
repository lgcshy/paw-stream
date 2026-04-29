<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Button } from 'vant'
import Layout from '@/components/Layout.vue'

const router = useRouter()

const goToStreams = () => {
  router.push('/streams')
}
</script>

<template>
  <Layout>
    <div class="home-view">
      <div class="content">
        <h2>{{ $t('home.welcome') }}</h2>
        <p>{{ $t('home.subtitle') }}</p>
        <Button type="primary" size="large" round @click="goToStreams"> {{ $t('home.viewStreams') }} </Button>
      </div>
    </div>
  </Layout>
</template>

<style scoped>
.home-view {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.content {
  text-align: center;
}

.content h2 {
  font-size: 24px;
  margin-bottom: 10px;
  color: #333;
}

.content p {
  font-size: 14px;
  color: #666;
  margin-bottom: 30px;
}
</style>
