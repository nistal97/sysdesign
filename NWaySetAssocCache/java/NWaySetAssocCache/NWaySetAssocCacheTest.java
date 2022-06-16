package NWaySetAssocCache;

import org.junit.Assert;
import org.junit.Test;

import java.util.concurrent.CountDownLatch;

public class NWaySetAssocCacheTest {

    @Test
    public void lru() throws Exception {
        putThenGet(Icache.CacheStrategy.LRU);
    }

    @Test
    public void mru() throws Exception {
        putThenGet(Icache.CacheStrategy.MRU);
    }

    void putThenGet(Icache.CacheStrategy strategy) throws Exception {
        Icache<Character, Integer> cache = Cache.getCache(strategy, 26, 26);
        for (char c = 'a'; c <= 'z'; c++) {
            cache.put(c, c - 'a');
        }
        for (char c = 'a'; c <= 'z'; c++) {
            Assert.assertTrue((int) cache.get(c) == c - 'a');
        }

        cache.put('A', 100);
        if (strategy == Icache.CacheStrategy.LRU) Assert.assertFalse(cache.contains('a'));
        if (strategy == Icache.CacheStrategy.MRU) Assert.assertFalse(cache.contains('z'));
    }

    @Test
    public void removeThenDump() throws Exception {
        Icache<Character, Integer> cache = Cache.getCache(Icache.CacheStrategy.LRU, 5, 2);
        for (char c = 'a'; c <= 'e'; c++) {
            cache.put(c, c - 'a');
        }
        cache.remove('a');
        Assert.assertFalse(cache.contains('a'));
        System.out.println(cache.dump());
    }

    @Test
    public void customEvicter() throws Exception {
        Exception ee = null;
        try {
            Icache<Character, Integer> cache = Cache.getCache(Icache.CacheStrategy.CUSTOM, 5, 2);
        } catch (Exception e) {
            ee = e;
        }
        Assert.assertTrue(ee != null);

        Icache<Character, Integer> cache = Cache.getCache(Icache.CacheStrategy.CUSTOM, 4, 4, () -> {
            return 'c';
        });
        for (char c = 'a'; c <= 'e'; c++) {
            cache.put(c, c - 'a');
        }
        Assert.assertFalse(cache.contains('c'));

    }

    @Test
    public void current() throws Exception {
        Icache<Integer, Integer> cache = Cache.getCache(Icache.CacheStrategy.LRU, 4, 4);
        CountDownLatch latch = new CountDownLatch(1000);
        for (int i = 0; i < 1000; i++) {
            final int s = i;
            Thread t = new Thread(() -> {
                cache.put(s, s);
                latch.countDown();
            });
            t.start();
        }
        latch.await();
        System.out.println(cache.dump());
    }

}