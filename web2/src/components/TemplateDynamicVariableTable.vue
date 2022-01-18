<template>
  <div>
    <span class="text-subtitle-1">Dynamic Variables</span>
    <v-row>
      <v-col cols="4" class="pb-0">
        Name
      </v-col>
      <v-col cols="4" class="pb-0">
        Default Value
      </v-col>
      <v-col cols="2" class="pb-0 d-flex justify-center">
        Required
      </v-col>
      <v-col cols="2" class="pb-0 d-flex justify-center">
        Remove
      </v-col>
    </v-row>
    <TemplateDynamicVariable
      v-for="(dynamicVar, i) in variables"
      :key="i"
      :variable="dynamicVar"
      :disabled="disabled"
      @remove="removeDynamicVariable(i)"
      @update="updateDynamicVariable($event, i)" />
    <v-row>
      <v-spacer></v-spacer>
      <v-col cols="3">
        <v-btn
            color="primary"
            @click="addDynamicVariable()"
        >New Variable</v-btn>
      </v-col>
    </v-row>
  </div>
</template>

<script>
/* eslint-disable import/extensions */
// TODO: add hint for "required"?
import TemplateDynamicVariable from '@/components/TemplateDynamicVariable';

export default {
  props: {
    variables: Array,
    disabled: Boolean,
  },

  components: {
    TemplateDynamicVariable,
  },

  methods: {
    addDynamicVariable() {
      const dynamicVariable = {
        name: '',
        default: '',
        required: true,
      };
      const newVariables = this.variables.concat([dynamicVariable]);
      this.$emit('change', newVariables);
    },

    removeDynamicVariable(index) {
      const newVariables = this.variables.slice();
      newVariables.splice(index, 1);
      this.$emit('change', newVariables);
    },

    updateDynamicVariable(updatedVariable, index) {
      const newVariables = this.variables.slice();
      newVariables.splice(index, 1, updatedVariable);
      this.$emit('change', newVariables);
    },
  },
};
</script>
