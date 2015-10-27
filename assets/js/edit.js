var viewModel = function(name, title, metaDescr, header, body, footer, key){
      var self = this;
      self.showSuccess = ko.observable(false);
      self.showFail = ko.observable(false);
      self.pageNameValue = ko.observable(name);
      self.pageTitleValue = ko.observable(title);
      self.pageMetaDescrValue = ko.observable(metaDescr);
      self.pageHeaderValue = ko.observable(header);
      self.pageBodyValue = ko.observable(body);
      self.pageFooterValue = ko.observable(footer);
      self.key = ko.observable(key);

      self.editSubmit = function(){
        self.post();
        return false;
      }

      self.post = function(){

        $.ajax({
            type: "POST",
            url: '/edit/post',
            data: 'pageNameValue=' + escape(self.pageNameValue()) + 
              '&pageTitleValue=' + escape(self.pageTitleValue()) + 
              '&pageMetaDescrValue=' + escape(self.pageMetaDescrValue()) + 
              '&pageHeaderValue=' + escape(self.pageHeaderValue()) + 
              '&pageBodyValue=' + escape(CKEDITOR.instances.editor1.getData()) + 
              '&pageFooterValue=' + escape(self.pageFooterValue()) + 
              '&key=' + escape(self.key()) + 
              '&pathname=' + location.pathname,
            success: function (data) {

                if(data.Name === 'Success' && data.Body === 'true') {
                  self.key(data.Key);
                  self.showSuccess(true);
                } 
                else {
                  self.showFail(true);
                }
            },
            error: function(obj, msg){
              self.showFail(true);
            },
            dataType: 'json',
            cache: false
        });
      }
    }