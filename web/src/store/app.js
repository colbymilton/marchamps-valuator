import { defineStore } from 'pinia'

export const useAppStore = () => {
  const innerStore = defineStore({
    "id": "app",
    state: () => ({
      packs: [],
      packValues: [],
      selectedPack: [],
      loaded: false,
    }),

    actions: {
      async getPacks() {
        this.loaded = true;
        const result = await fetch('http://localhost:9999/packs');
        const data = await result.json();
        this.packs = data;

        for (let i = 0; i < this.packs.length; i++) {
          let pack = this.packs[i];
          if (pack.code == "core") {
            pack.owned = true;
            pack.locked = true;
          }
        }
      },

      getPack(code) {
        for (let i = 0; i < this.packs.length; i++) {
          let pack = this.packs[i];
          if (pack.code == code) {
            return pack;
          }
        }
        return [];
      },

      packsString() {
        let s = "";
        for (let i = 0; i < this.packs.length; i++) {
          let pack = this.packs[i];
          if (pack.owned) {
            s += pack.code + ","
          }
        }
        return s;
      }
    },
  });

  const s = innerStore();
  if (!s.loaded) {
    s.getPacks();
  }
  return s;
}
