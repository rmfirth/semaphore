<template>
  <v-form
      ref="form"
      lazy-validation
      v-model="formValid"
      v-if="isLoaded()"
  >
    <v-alert
        :value="formError"
        color="error"
        class="pb-2"
    >{{ formError }}
    </v-alert>

    <v-alert
        color="blue"
        dark
        icon="mdi-source-fork"
        dismissible
        v-model="commitAvailable"
        prominent
    >
      <div
          style="font-weight: bold;"
      >{{ commitHash ? commitHash.substr(0, 10) : '' }}
      </div>
      <div v-if="commitMessage">{{ commitMessage }}</div>
    </v-alert>

    <v-select
        v-if="template.type === 'deploy'"
        v-model="item.build_task_id"
        label="Build Version"
        :items="buildTasks"
        item-value="id"
        :item-text="(itm) => itm.version + (itm.message ? ' — ' + itm.message : '')"
        :rules="[v => !!v || 'Build Version is required']"
        required
        :disabled="formSaving"
    />

    <v-text-field
        v-model="item.message"
        label="Message (Optional)"
        :disabled="formSaving"
    />

    <v-row no-gutters>
      <v-col>
        <v-checkbox
            v-model="item.debug"
            label="Debug"
        ></v-checkbox>
      </v-col>
      <v-col>
        <v-checkbox
            v-model="item.dry_run"
            label="Dry Run"
        ></v-checkbox>
      </v-col>
    </v-row>

    <v-row v-if="dynamicVars.length > 0">
      <v-col>
        <span class="text-subtitle-1 d-block mb-4">Dynamic Variables</span>
        <v-text-field
          v-for="(variable, index) in dynamicVars"
          :value ="variable.value"
          @input="updateDynamicVariables($event, index)"
          :key="variable.name"
          :label="variable.name"
          :rules="[v => (!variable.required || !!v) || `${variable.name} is required`]"
          :disabled="formSaving"
          class="py-0 my-4"
        ></v-text-field>
      </v-col>
    </v-row>

  </v-form>
</template>
<script>
import ItemFormBase from '@/components/ItemFormBase';
import axios from 'axios';

export default {
  mixins: [ItemFormBase],
  props: {
    templateId: Number,
    commitHash: String,
    commitMessage: String,
    buildTask: Object,
  },
  data() {
    return {
      template: null,
      buildTasks: null,
      commitAvailable: null,
      dynamicVars: null,
    };
  },

  watch: {
    needReset(val) {
      if (val) {
        this.item.template_id = this.templateId;
      }
    },

    templateId(val) {
      this.item.template_id = val;
    },

    commitHash(val) {
      this.item.commit_hash = val;
      this.commitAvailable = this.item.commit_hash != null;
    },

    version(val) {
      this.item.version = val;
    },

    commitAvailable(val) {
      this.item.commit_hash = val ? this.commitHash : null;
    },
  },

  methods: {
    isLoaded() {
      return this.item != null
          && this.template != null
          && this.buildTasks != null
          && this.dynamicVars != null;
    },

    async afterLoadData() {
      this.item.template_id = this.templateId;

      this.template = (await axios({
        keys: 'get',
        url: `/api/project/${this.projectId}/templates/${this.templateId}`,
        responseType: 'json',
      })).data;

      if (!(this.template.dynamic_vars === '' || this.template.dynamic_vars === undefined)) {
        this.dynamicVars = JSON.parse(this.template.dynamic_vars);
      } else {
        this.dynamicVars = [];
      }

      this.buildTasks = this.template.type === 'deploy' ? (await axios({
        keys: 'get',
        url: `/api/project/${this.projectId}/templates/${this.template.build_template_id}/tasks`,
        responseType: 'json',
      })).data.filter((task) => task.version != null && task.status === 'success') : [];

      if (this.buildTasks.length > 0) {
        this.item.build_task_id = this.build_task ? this.build_task.id : this.buildTasks[0].id;
      }

      this.commitAvailable = this.commitHash != null;
    },

    getItemsUrl() {
      return `/api/project/${this.projectId}/tasks`;
    },

    updateDynamicVariables(variableValue, index) {
      const variables = this.dynamicVars;
      const newVariable = variables[index];
      newVariable.value = variableValue;
      const newDynamicVars = this.dynamicVars.slice();
      newDynamicVars.splice(index, 1, newVariable);
      this.dynamicVars = newDynamicVars;
    },

    async beforeSave() {
      const dynamicVars = {};
      this.dynamicVars.forEach((item) => {
        const itemKey = item.name;
        const itemValue = item.value;
        dynamicVars[itemKey] = itemValue;
      });
      this.item.dynamic_vars = JSON.stringify(dynamicVars);
    },
  },
};
</script>
