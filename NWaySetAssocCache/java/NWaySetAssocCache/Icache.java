package NWaySetAssocCache;

public interface Icache<K, V> {
    V get(K k);
    void put(K k, V v);
    boolean contains(K k);
    int size();
    void remove(K k);
    StringBuilder dump();

    enum CacheStrategy {
        LRU,
        MRU,
        CUSTOM
    }
}
