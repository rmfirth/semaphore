<template>
  <v-row align="center">
    <v-col cols="4" class="py-0">
      <v-text-field
          :value="variable.name"
          @input="updateName($event)"
          :rules="[v => !!v || 'Variable name is required']"
          :disabled="disabled"
      ></v-text-field>
    </v-col>
    <v-col cols="4" class="py-0">
      <v-text-field
          :value="variable.default"
          @input="updateDefaultValue($event)"
          :disabled="defaultValueDisabled"
      ></v-text-field>
    </v-col>
    <v-col cols="2" class="py-0 d-flex justify-center">
        <v-switch
          :input-value="variable.required"
          @change="updateRequired($event)"
          :disabled="disabled" />
    </v-col>
    <v-col cols="2" class="py-0 d-flex justify-center">
        <v-btn
          icon
          color="error"
          @click="$emit('remove')"
          :disabled="disabled"
        >
          <v-icon>mdi-delete</v-icon>
        </v-btn>
    </v-col>
  </v-row>
</template>

<script>

export default {
  props: {
    variable: Object,
    disabled: Boolean,
  },

  computed: {
    defaultValueDisabled() {
      return this.disabled || this.variable.required;
    },
  },

  methods: {
    clearDefaultValue() {
      this.updateDefaultValue('');
    },

    updateName(newName) {
      const updatedVariable = { ...this.variable };
      updatedVariable.name = newName;
      this.$emit('update', updatedVariable);
    },

    updateDefaultValue(newDefault) {
      const updatedVariable = { ...this.variable };
      updatedVariable.default = newDefault;
      this.$emit('update', updatedVariable);
    },

    updateRequired(newValue) {
      const updatedVariable = { ...this.variable };
      updatedVariable.required = newValue;
      this.$emit('update', updatedVariable);
    },
  },
};
</script>
