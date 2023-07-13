import { defineStore } from 'pinia'

export const useAppStore = defineStore({
  "id": "app",
  state: () => ({
    packs: [],
    packValues: [],
    selectedPack: [],
    error: "",
  }),

  actions: {
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
          s += pack.code + ",";
        }
      }
      return s;
    }
  },
});
