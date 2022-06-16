package NWaySetAssocCache;

import java.util.*;
import java.util.concurrent.locks.ReentrantLock;

public class Cache<K, V> {

    public static <K, V>Icache getCache(Icache.CacheStrategy strategy, int cap, int limit, IEvicter... evicter) throws Exception {
        if (cap < limit) throw new Exception("cap should be bigger than limit");
        if (limit <= 0) throw new Exception("limit should > 0");

        if (strategy == Icache.CacheStrategy.CUSTOM && evicter.length == 0) {
            throw new Exception("No custom evicter provided");
        }

        return new CacheImpl(strategy, cap, limit, evicter.length > 0 ? evicter[0] : null);
    }

    private static class CacheImpl<K, V> implements Icache<K, V> {
        private Icache.CacheStrategy strategy;
        private IEvicter evicter;
        private int cap;
        private int limit;
        private int sets;
        private int size = 0;
        private List<Map<K, AddrV>> maps = new ArrayList<>();
        private List<ListNode> list = new ArrayList<>();

        public CacheImpl(Icache.CacheStrategy strategy, int cap, int limit, IEvicter evicter) {
            this.strategy = strategy;
            this.evicter = evicter;
            this.cap = cap;
            this.limit = limit;
            this.sets = (int)Math.ceil((cap/limit));
            for (int i = 0;i < sets;i ++) {
                maps.add(new HashMap());
                list.add(new ListNode());
            }
        }

        class AddrV<K, V> {
            K k;
            V v;
            AddrV prev;
            AddrV next;
            AddrV(K k, V v) { this.k = k; this.v = v;}
        }
        class ListNode {
            AddrV head;
            AddrV tail;
            int len;
            ReentrantLock lck = new ReentrantLock();
        }

        @Override
        public V get(K k) {
            int mod = locate(k);
            try {
                list.get(mod).lck.lock();
                if (maps.get(mod).containsKey(k)) {
                    AddrV v = maps.get(mod).get(k);
                    Objects.requireNonNull(v);
                    remove(list.get(mod), v);
                    pushFront(list.get(mod), v);
                    return (V) maps.get(mod).get(k).v;
                }
            } finally {
                list.get(mod).lck.unlock();
            }
            return null;
        }

        @Override
        public void put(K k, V v) {
            int mod = locate(k);
            try {
                list.get(mod).lck.lock();
                if (list.get(mod).len == limit) {
                    AddrV t = null;
                    if (strategy == CacheStrategy.LRU) {
                        t = list.get(mod).tail;
                    } else if (strategy == CacheStrategy.MRU) {
                        t = list.get(mod).head;
                    } else {
                        t = maps.get(mod).get((evicter.evict()));
                    }
                    if (t != null) {
                        remove(list.get(mod), t);
                        maps.get(mod).remove(t.k);
                        size --;
                    }
                }

                AddrV n = new AddrV(k, v);
                pushFront(list.get(mod), n);
                maps.get(mod).put(k, n);
                size ++;
            } finally {
                list.get(mod).lck.unlock();
            }
        }

        @Override
        public boolean contains(K k) {
            int mod = locate(k);
            try {
                list.get(mod).lck.lock();
                if (maps.get(mod).containsKey(k)) return true;
            } finally {
                list.get(mod).lck.unlock();
            }
            return false;
        }

        @Override
        public int size() {
            return size;
        }

        @Override
        public void remove(K k) {
            int mod = locate(k);
            try {
                list.get(mod).lck.lock();
                if (contains(k)) {
                    AddrV v = maps.get(mod).get(k);
                    remove(list.get(mod), v);
                    maps.get(mod).remove(k);
                    size --;
                }
            } finally {
                list.get(mod).lck.unlock();
            }
        }

        @Override
        public StringBuilder dump() {
            StringBuilder sb = new StringBuilder();
            sb.append("cap:").append(cap).append(" limit:").append(limit).append(" size:").append(size);
            for (int i = 0;i < sets;i ++) {
                AddrV t = list.get(i).head;
                sb.append("[");
                while (t != null) {
                    sb.append("{K:").append(t.k).append(" V:").append(t.v).append("}");
                    t = t.next;
                }
                sb.append("]");
            }
            sb.append("\n");
            return sb;
        }

        private int locate(K k) {
            Objects.requireNonNull(k);
            return k.hashCode() % sets;
        }

        private AddrV pushFront(ListNode node, AddrV v) {
            if (v == null) return v;

            if (node.head != null) {
                v.next = node.head;
                node.head.prev = v;
            } else {
                node.tail = v;
            }
            node.head = v;
            node.len ++;
            return v;
        }

        private AddrV remove(ListNode node, AddrV v) {
            if (v == null) return v;

            if (v.prev != null && v.next != null) {
                v.prev.next = v.next;
                v.next.prev = v.prev;
            } else if (v.prev == null && v.next == null) {
                node.head = v;
                node.tail = v;
            } else {
                if (v.prev == null) {
                    node.head = v.next;
                    v.next.prev = null;
                } else {
                    node.tail = v.prev;
                    v.prev.next = null;
                }
            }
            node.len --;
            return v;
        }
    }

}
